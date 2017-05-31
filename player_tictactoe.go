package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/boardgamesai/games/tictactoe"
)

func main() {
	// First thing we do upon launch is let our invoker know we started up okay.
	// There could be Go compile-time issues preventing us from getting here.
	fmt.Printf("OK\n")

	stdin := bufio.NewReader(os.Stdin)

	for {
		// Get raw JSON
		inputStr, err := stdin.ReadString('\n')
		if err != nil && err != io.EOF {
			log.Fatalf("Error reading input: %s\n", err)
		}

		// Convert JSON to Message
		message := tictactoe.Message{}
		err = json.Unmarshal([]byte(inputStr), &message)
		if err != nil {
			log.Fatalf("Error decoding JSON: %s err: %s", inputStr, err)
		}
		board := tictactoe.GetBoardFromString(message.Board)

		// Get Move - needs to be defined in a file next to this one
		move := GetMove(message.Symbol, board)

		// Convert Move to JSON
		moveJson, err := json.Marshal(&move)
		if err != nil {
			log.Fatalf("Couldn't convert move to JSON: %+v", move)
		}

		// Write JSON back and we're done!
		fmt.Printf("%s\n", moveJson)
	}
}
