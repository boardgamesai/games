package game

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime/debug"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

const (
	PlayerLaunchTimeout   = 30
	PlayerResponseTimeout = 15
)

type RunnablePlayer struct {
	gameName     string
	filePath     string // Path of stored user-written code
	runDir       string // The tmp dir where this player is running
	cmd          *exec.Cmd
	cmdStdin     *io.WriteCloser
	cmdStdout    *bufio.Reader
	cmdStderr    *bytes.Buffer
	responseChan chan []byte
}

func NewRunnablePlayer(gameName string, filePath string) *RunnablePlayer {
	player := RunnablePlayer{
		gameName: gameName,
		filePath: filePath,
	}
	return &player
}

func (p *RunnablePlayer) Run() error {
	if err := p.setupFiles(); err != nil {
		return err
	}

	if err := p.launchProcess(); err != nil {
		return err
	}

	// Spin off a goroutine to read the "OK" response so we can block on it below.
	go p.readResponseAsync()

	var err error
	select {
	case response := <-p.responseChan:
		if string(response) != "OK" {
			err = fmt.Errorf("Got non-OK response when launching player: %s stderr: %s", response, p.Stderr())
		}
	case <-time.After(time.Second * PlayerLaunchTimeout):
		err = fmt.Errorf("Timeout launching player")
	}

	return err
}

func (p *RunnablePlayer) Build() (string, error) {
	if err := p.setupFiles(); err != nil {
		return "", err
	}
	defer p.CleanUp()

	compiledOutput := p.runDir + "/ai"
	cmd := exec.Command("go", "build", "-o", compiledOutput, p.runDir+"/main.go", p.runDir+"/ai.go")
	outputBytes, _ := cmd.CombinedOutput()
	output := string(outputBytes)

	// Last thing, remove any references to the runDir from the error string
	output = strings.ReplaceAll(output, p.runDir+"/", "")

	return string(output), nil
}

func (p *RunnablePlayer) CleanUp() error {
	// First wipe out the tmp dir where we copied everything.
	err1 := os.RemoveAll(p.runDir)

	// Kill the process (if it didn't die already due to error).
	var err2 error
	if p.cmd != nil {
		err2 = p.cmd.Process.Kill()
	}

	if err1 != nil {
		return err1
	}
	return err2
}

func (p *RunnablePlayer) SendMessage(message interface{}) ([]byte, error) {
	// Let's use reflection to get the type of this message
	messageType := reflect.TypeOf(message).Name()
	if messageType[0:7] != "Message" {
		return []byte{}, fmt.Errorf("Invalid type %s passed to SendMessage", messageType)
	}

	// Hack off the "Message" on the front and lowercase it
	messageType = strings.ToLower(messageType[7:])

	messageJSON, err := json.Marshal(&message)
	if err != nil {
		return []byte{}, err
	}

	m := Message{
		Type: messageType,
		Data: messageJSON,
	}

	mJSON, err := json.Marshal(&m)
	if err != nil {
		return []byte{}, err
	}

	err = p.writeLine(string(mJSON))
	if err != nil {
		return []byte{}, err
	}

	// This spins off a goroutine to read the response, and we block on it just below.
	go p.readResponseAsync()

	var response []byte

	select {
	case response = <-p.responseChan:
		// Do nothing, the assignment above is the important thing
	case <-time.After(time.Second * PlayerResponseTimeout):
		err = fmt.Errorf("Timeout reading player response")
	}

	return response, err
}

func (p *RunnablePlayer) SendMessageNoResponse(message interface{}) error {
	response, err := p.SendMessage(message)
	if err != nil {
		return err
	}
	if string(response) != "OK" {
		return fmt.Errorf("Got non-OK response: %s stderr: %s", response, p.Stderr())
	}

	return nil
}

func (p *RunnablePlayer) Stderr() string {
	if p.cmdStderr == nil {
		return ""
	}

	buf := *p.cmdStderr
	return buf.String()
}

func (p *RunnablePlayer) String() string {
	return p.filePath
}

func (p *RunnablePlayer) setupFiles() error {
	config, err := Config()
	if err != nil {
		return err
	}

	// Ensure this player exists
	aiSrcPath := p.filePath
	if _, err := os.Stat(aiSrcPath); os.IsNotExist(err) {
		return fmt.Errorf("Player file does not exist: %s", aiSrcPath)
	}

	// First create the tmp dir for the player
	tmpDir := config.TmpDir + "/" + uuid.NewRandom().String()
	if err := os.Mkdir(tmpDir, 0700); err != nil {
		return fmt.Errorf("Could not create tmp dir: %s for player: %s err: %s", tmpDir, p.filePath, err)
	}
	p.runDir = tmpDir

	// Next copy over the main driver file
	srcPath, err := p.driverFilePath()
	if err != nil {
		return err
	}
	destPath := tmpDir + "/main.go"
	err = p.copyFile(srcPath, destPath)
	if err != nil {
		return fmt.Errorf("Could not copy %s to %s", srcPath, destPath)
	}

	// Now copy over the AI-specific file
	aiDestPath := tmpDir + "/ai.go"
	err = p.copyFile(aiSrcPath, aiDestPath)
	if err != nil {
		return fmt.Errorf("Could not copy %s to %s", aiSrcPath, aiDestPath)
	}

	return nil
}

func (p *RunnablePlayer) driverFilePath() (string, error) {
	// First check if we've cloned the games repo and are developing against it directly.
	path := p.gameName + "/ai/main.go"
	_, err := os.Stat(path)
	if err == nil {
		return path, nil
	} else if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	// This means we were brought in via go modules, so we need to find which version is used
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return "", errors.New("not using Go modules and cannot locate games lib")
	}
	version := ""
	for _, dep := range buildInfo.Deps {
		if strings.Contains("boardgamesai/games", dep.Path) {
			path = dep.Path
			version = dep.Version
			break
		}
	}

	if version == "" {
		return "", errors.New("could not determine Go module version of games lib")
	}

	return os.Getenv("GOPATH") + "/pkg/mod/" + path + "@" + version + "/" + p.gameName + "/ai/main.go", nil
}

func (p *RunnablePlayer) launchProcess() error {
	cmd := exec.Command("go", "run", p.runDir+"/main.go", p.runDir+"/ai.go")
	p.cmd = cmd

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	p.cmdStdin = &stdin

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	p.cmdStdout = bufio.NewReader(stdout)

	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	p.cmdStderr = &errBuf

	p.responseChan = make(chan []byte, 1)

	return cmd.Start()
}

func (p *RunnablePlayer) writeLine(line string) error {
	_, err := io.WriteString(*p.cmdStdin, fmt.Sprintf("%s\n", line))
	return err
}

func (p *RunnablePlayer) readResponseAsync() {
	response, err := p.cmdStdout.ReadBytes('\n')
	if err != nil && err != io.EOF {
		// TODO how to communicate this error?
		log.Fatalf("couldn't read from stdin: %s\n", err)
		return
	}

	// Chop off the newline, if there is one.
	responseLen := len(response)
	if responseLen >= 1 && response[responseLen-1] == '\n' {
		response = response[:responseLen-1]
	}

	p.responseChan <- response
}

func (p *RunnablePlayer) copyFile(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func FileNameToPlayerName(filePath string) string {
	filename := filepath.Base(filePath)
	if len(filename) >= 3 && filename[len(filename)-3:] == ".go" {
		filename = filename[0 : len(filename)-3]
	}
	return filename
}
