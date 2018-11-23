package game

import (
	"bufio"
	"encoding/json"
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
	mJSON, err := d.stdin.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return Message{}, err
	}

	m := Message{}
	err = json.Unmarshal(mJSON, &m)
	return m, err
}
