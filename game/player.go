package game

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"

	"github.com/boardgamesai/games/util"
	"github.com/pborman/uuid"
)

const (
	PlayerLaunchTimeout   = 30
	PlayerResponseTimeout = 15
)

type Player struct {
	gameName     string
	fileName     string // Non-absolute filename of stored user-written code
	runDir       string // The tmp dir where this player is running
	cmd          *exec.Cmd
	cmdStdin     *io.WriteCloser
	cmdStdout    *bufio.Reader
	cmdStderr    *bytes.Buffer
	responseChan chan []byte
}

func NewPlayer(gameName string, playerName string) *Player {
	player := Player{
		gameName: gameName,
		fileName: playerName,
	}
	return &player
}

func (p *Player) setupFiles(config *Configuration) error {
	// Ensure this player exists
	aiSrcPath := os.Getenv("GOPATH") + config.PlayerDir + "/" + p.gameName + "/" + p.fileName + "/" + p.fileName + ".go"
	if _, err := os.Stat(aiSrcPath); os.IsNotExist(err) {
		return fmt.Errorf("Player file does not exist: %s", aiSrcPath)
	}

	// First create the tmp dir for the player
	tmpDir := os.Getenv("GOPATH") + config.TmpDir + "/" + uuid.NewRandom().String()
	err := os.Mkdir(tmpDir, 0700)
	if err != nil {
		return fmt.Errorf("Could not create tmp dir: %s for player: %s err: %s", tmpDir, p.fileName, err)
	}
	p.runDir = tmpDir

	// Next copy over the main driver file
	srcPath := p.gameName + "/ai/main.go"
	destPath := tmpDir + "/main.go"
	err = util.CopyFile(srcPath, destPath)
	if err != nil {
		return fmt.Errorf("Could not copy %s to %s", srcPath, destPath)
	}

	// Now copy over the AI-specific file
	aiDestPath := tmpDir + "/ai.go"
	err = util.CopyFile(aiSrcPath, aiDestPath)
	if err != nil {
		return fmt.Errorf("Could not copy %s to %s", aiSrcPath, aiDestPath)
	}

	return nil
}

func (p *Player) launchProcess(config *Configuration) error {
	cmd := exec.Command("go", "run", p.runDir+"/main.go", p.runDir+"/ai.go")
	if config.UseSandbox {
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "GOOS=nacl")
		cmd.Env = append(cmd.Env, "GOARCH=amd64p32")
	}
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

func (p *Player) Run(config *Configuration) error {
	err := p.setupFiles(config)
	if err != nil {
		return err
	}

	err = p.launchProcess(config)
	if err != nil {
		return err
	}

	// Spin off a goroutine to read the "OK" response so we can block on it below.
	go p.readResponseAsync()

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

func (p *Player) CleanUp() error {
	// First wipe out the tmp dir where we copied everything.
	err1 := os.RemoveAll(p.runDir)

	// Kill the process (if it didn't die already due to error).
	err2 := p.cmd.Process.Kill()

	if err1 != nil {
		return err1
	}
	return err2
}

// Our messages have two lines. First line is the type, second line is the JSON-encoded payload.
func (p *Player) SendMessage(data interface{}) ([]byte, error) {
	// Let's use reflection to get the type of this message
	messageType := reflect.TypeOf(data).Name()
	// Hack off the "Message" on the front and lowercase it
	messageType = strings.ToLower(messageType[7:])

	dataJSON, err := json.Marshal(&data)
	if err != nil {
		return []byte{}, err
	}

	err = p.writeLine(messageType)
	if err != nil {
		return []byte{}, err
	}

	err = p.writeLine(string(dataJSON))
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

func (p *Player) writeLine(line string) error {
	_, err := io.WriteString(*p.cmdStdin, fmt.Sprintf("%s\n", line))
	return err
}

func (p *Player) readResponseAsync() {
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

func (p *Player) Stderr() string {
	buf := *p.cmdStderr
	return buf.String()
}

func (p *Player) String() string {
	return p.fileName
}
