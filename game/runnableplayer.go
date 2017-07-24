package game

import (
	"bufio"
	"bytes"
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

type RunnablePlayer struct {
	PlayerPath   string // The boardgamesAI-provided driver
	AIPath       string // The code the user writes
	cmd          *exec.Cmd
	cmdStdin     *io.WriteCloser
	cmdStdout    *bufio.Reader
	cmdStderr    *bytes.Buffer
	responseChan chan string
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

	// Wait here to make sure things started okay.
	p.readResponseAsync()

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

func (p *RunnablePlayer) SendMessage(messageJSON string) error {
	_, err := io.WriteString(*p.cmdStdin, fmt.Sprintf("%s\n", messageJSON))
	return err
}

func (p *RunnablePlayer) ReadResponse() (string, error) {
	p.readResponseAsync()

	var responseJSON string
	var err error

	select {
	case responseJSON = <-p.responseChan:
		// Do nothing, the assignment above is the important thing
	case <-time.After(time.Second * PlayerMoveTimeout):
		err = fmt.Errorf("Timeout reading player move response")
	}

	return responseJSON, err
}

func (p *RunnablePlayer) readResponseAsync() {
	go func() {
		response, err := p.cmdStdout.ReadString('\n')
		if err != nil && err != io.EOF {
			// TODO how to communicate this error?
			log.Fatalf("couldn't read from stdin: %s\n", err)
		}

		p.responseChan <- strings.TrimSpace(response)
	}()
}

func (p *RunnablePlayer) Stderr() string {
	buf := *p.cmdStderr
	return buf.String()
}

func (p *RunnablePlayer) String() string {
	return p.AIPath
}
