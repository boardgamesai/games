package game

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/boardgamesai/games/util"
	"github.com/pborman/uuid"
)

const (
	PlayerLaunchTimeout = 30
	PlayerMoveTimeout   = 15
)

type RunnablePlayer struct {
	PlayerPath   string // The boardgamesAI-provided driver
	AIPath       string // The code the user writes
	cmd          *exec.Cmd
	cmdStdin     *io.WriteCloser
	cmdStdout    *bufio.Reader
	cmdStderr    *bytes.Buffer
	responseChan chan string
}

func NewRunnablePlayer(config *Configuration, gameName string, playerName string) (*RunnablePlayer, error) {
	// Ensure this player exists
	aiSrcPath := os.Getenv("GOPATH") + config.PlayerDir + "/" + gameName + "/" + playerName + "/" + playerName + ".go"
	if _, err := os.Stat(aiSrcPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("Player file does not exist: %s", aiSrcPath)
	}

	// First create the tmp dir for the player
	tmpDir := os.Getenv("GOPATH") + config.TmpDir + "/" + uuid.NewRandom().String()
	err := os.Mkdir(tmpDir, 0700)
	if err != nil {
		return nil, fmt.Errorf("Could not create tmp dir: %s for player: %s err: %s", tmpDir, playerName, err)
	}

	// Next copy over the base player file
	playerFile := "player_" + gameName + ".go"
	playerDestPath := tmpDir + "/" + playerFile
	err = util.CopyFile(playerFile, playerDestPath)
	if err != nil {
		return nil, fmt.Errorf("Could not copy %s to %s", playerFile, playerDestPath)
	}

	// Now copy over the AI-specific file
	aiDestPath := tmpDir + "/" + playerName + ".go"
	err = util.CopyFile(aiSrcPath, aiDestPath)
	if err != nil {
		return nil, fmt.Errorf("Could not copy %s to %s", aiSrcPath, aiDestPath)
	}

	player := RunnablePlayer{
		PlayerPath: playerDestPath,
		AIPath:     aiDestPath,
	}
	return &player, nil
}

func (p *RunnablePlayer) Run(useSandbox bool) error {
	cmd := exec.Command("go", "run", p.PlayerPath, p.AIPath)
	if useSandbox {
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

	p.responseChan = make(chan string, 1)

	if err = cmd.Start(); err != nil {
		return err
	}

	// Spin off a goroutine to read the "OK" response so we can block on it below.
	go p.readResponseAsync()

	select {
	case response := <-p.responseChan:
		if response != "OK" {
			err = fmt.Errorf("Got non-OK response when launching player: %s stderr: %s", response, p.Stderr())
		}
	case <-time.After(time.Second * PlayerLaunchTimeout):
		err = fmt.Errorf("Timeout launching player")
	}

	return err
}

func (p *RunnablePlayer) CleanUp() error {
	return os.RemoveAll(filepath.Dir(p.PlayerPath))
}

func (p *RunnablePlayer) SendMessage(messageJSON string) (string, error) {
	_, err := io.WriteString(*p.cmdStdin, fmt.Sprintf("%s\n", messageJSON))
	if err != nil {
		return "", err
	}

	// This spins off a goroutine to read the response, and we block on it just below.
	go p.readResponseAsync()

	var response string

	select {
	case response = <-p.responseChan:
		// Do nothing, the assignment above is the important thing
	case <-time.After(time.Second * PlayerMoveTimeout):
		err = fmt.Errorf("Timeout reading player move response")
	}

	return response, err
}

func (p *RunnablePlayer) readResponseAsync() {
	response, err := p.cmdStdout.ReadString('\n')
	if err != nil && err != io.EOF {
		// TODO how to communicate this error?
		log.Fatalf("couldn't read from stdin: %s\n", err)
	}

	p.responseChan <- strings.TrimSpace(response)
}

func (p *RunnablePlayer) Stderr() string {
	buf := *p.cmdStderr
	return buf.String()
}

func (p *RunnablePlayer) String() string {
	return p.AIPath
}
