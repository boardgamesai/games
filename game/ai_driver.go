package game

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type AIDriver struct {
	stdin *bufio.Reader
}

func (d *AIDriver) Setup() {
	// First thing we do upon launch is let our invoker know we started up okay.
	// There could be Go compile-time issues preventing us from getting here.
	fmt.Println("OK")

	// Now grab stdio, we need it for reading input later.
	d.stdin = bufio.NewReader(os.Stdin)
}

func (d *AIDriver) GetNextMessage() (Message, error) {
	// First line is the type
	messageType, err := d.stdin.ReadString('\n')
	if err != nil && err != io.EOF {
		return Message{}, err
	}
	messageType = messageType[:len(messageType)-1] // Remove trailing newline, ReadString() includes it

	// Second line is the JSON payload
	messageJSON, err := d.stdin.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return Message{}, err
	}

	message := Message{
		Type: messageType,
		Data: messageJSON,
	}
	return message, nil
}
