package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/boardgamesai/games/fourinarow"
)

func main() {
	// First thing we do upon launch is let our invoker know we started up okay.
	// There could be Go compile-time issues preventing us from getting here.
	fmt.Printf("OK\n")

	stdin := bufio.NewReader(os.Stdin)

	// Next thing after startup is to wait to be told our initial state.
	messageJSON, err := ReadLine(stdin)
	if err != nil {
		log.Fatalf("Error reading input: %s\n", err)
	}
	message := fourinarow.MessageSetup{}
	err = json.Unmarshal([]byte(messageJSON), &message)
	if err != nil {
		log.Fatalf("Error decoding JSON: %s err: %s", messageJSON, err)
	}
	fmt.Printf("OK\n")

	state := fourinarow.State{
		Order:   message.Order,
		Players: []fourinarow.Player{},
	}
	for _, player := range message.Players {
		player.PlayerPath = "" // Wipe out irrelevant stuff
		player.AIPath = ""     // TODO - don't send it in the first place
		state.Players = append(state.Players, *player)
	}

	for {
		// Get raw JSON
		messageJSON, err = ReadLine(stdin)
		if err != nil {
			log.Fatalf("Error reading input: %s\n", err)
		}

		// Convert JSON to Message
		message := fourinarow.MessageMove{}
		err = json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			log.Fatalf("Error decoding JSON: %s err: %s", messageJSON, err)
		}

		state.Board = fourinarow.GetBoardFromString(message.Board)
		state.NewMoves = message.NewMoves
		state.AllMoves = append(state.AllMoves, message.NewMoves...)

		// Get Move - needs to be defined in a file next to this one
		move := GetMove(&state)

		// Add new move to our state immediately - we don't get our own moves in NewMoves
		moveLog := fourinarow.MoveLog{
			Move:  move,
			Order: state.Order,
		}
		state.AllMoves = append(state.AllMoves, moveLog)

		// Convert Move to JSON
		moveJSON, err := json.Marshal(&move)
		if err != nil {
			log.Fatalf("Couldn't convert move to JSON: %+v", move)
		}

		// Write JSON back and we're done!
		fmt.Printf("%s\n", moveJSON)
	}
}

// ReadLine blocks until it gets something to read
func ReadLine(stdin *bufio.Reader) (string, error) {
	line, err := stdin.ReadString('\n')
	if err != nil && err != io.EOF {
		return "", err
	}

	return line, nil
}
