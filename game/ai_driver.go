package game

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

func (d *AIDriver) PrintResponse(data []byte) error {
	m := MessageResponse{
		Data: data,
	}

	return d.doPrint(m)
}

func (d *AIDriver) PrintErrorResponse(err *DQError) error {
	m := MessageResponse{
		Err: err,
	}
	return d.doPrint(m)
}

func (d *AIDriver) doPrint(m MessageResponse) error {
	messageJSON, err := json.Marshal(m)
	if err != nil {
		log.Printf("error printing response: %s", err)
		return err
	}

	// Need to strip any newlines since we use them to denote EOF when reading
	fmt.Println(strings.ReplaceAll(string(messageJSON), "\n", " "))
	return nil
}

func (d *AIDriver) OkJSON() []byte {
	strJSON, _ := json.Marshal("OK")
	return strJSON
}
