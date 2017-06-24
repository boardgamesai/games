package tictactoe

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	PlayerLaunchTimeout = 30
	PlayerMoveTimeout   = 15
)

type Player struct {
	ID         string
	Name       string
	Symbol     string // "X" or "O"
	Order      int    // 1-based
	PlayerPath string // The boardgamesAI-provided driver
	AIPath     string // The code the user writes
	cmd        *exec.Cmd
	cmdStdin   *io.WriteCloser
	cmdStdout  *bufio.Reader
	cmdStderr  *bytes.Buffer
	ch         chan string
}

func (p *Player) Run(useSandbox bool) error {
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

	p.ch = make(chan string, 1)

	if err = cmd.Start(); err != nil {
		return err
	}

	// Wait here to make sure things started okay.
	p.readResponseAsync()

	select {
	case response := <-p.ch:
		if response != "OK" {
			err = fmt.Errorf("Got non-OK response when launching player: %s stderr: %s", response, p.Stderr())
		}
	case <-time.After(time.Second * PlayerLaunchTimeout):
		err = fmt.Errorf("Timeout launching player")
	}

	if err != nil {
		return err
	}

	// Now tell them which symbol they are. Note that we don't wait for a response.
	// The other end reads line-by-line and will do the right thing.
	message := MessageSetup{
		Symbol: p.Symbol,
	}
	messageJSON, err := json.Marshal(&message)
	if err != nil {
		return err
	}
	_, err = io.WriteString(*p.cmdStdin, fmt.Sprintf("%s\n", messageJSON))
	if err != nil {
		return err
	}

	return nil
}

func (p *Player) GetMove(b *Board) (Move, error) {
	move := Move{}

	message := MessageMove{
		Board: GetStringFromBoard(b),
	}

	messageJSON, err := json.Marshal(&message)
	if err != nil {
		return move, err
	}

	// Write the state of things to the player
	_, err = io.WriteString(*p.cmdStdin, fmt.Sprintf("%s\n", messageJSON))
	if err != nil {
		return move, err
	}

	p.readResponseAsync()

	select {
	case moveJSON := <-p.ch:
		err = json.Unmarshal([]byte(moveJSON), &move)
	case <-time.After(time.Second * PlayerMoveTimeout):
		err = fmt.Errorf("Timeout getting player move")
	}

	return move, err
}

func (p *Player) readResponseAsync() {
	go func() {
		response, err := p.cmdStdout.ReadString('\n')
		if err != nil && err != io.EOF {
			// TODO how to communicate this error?
			log.Fatalf("couldn't read from stdin: %s\n", err)
		}

		p.ch <- strings.TrimSpace(response)
	}()
}

func (p *Player) Stderr() string {
	buf := *p.cmdStderr
	return buf.String()
}

func (p *Player) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Symbol)
}
