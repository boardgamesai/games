package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/tictactoe"
)

type AIDriver struct {
	game.AIDriver
	state *State
	ai    tictactoeAI
}

func New(ai tictactoeAI) *AIDriver {
	return &AIDriver{
		state: &State{},
		ai:    ai,
	}
}

func (d *AIDriver) Run() {
	d.Setup()

	for {
		message, err := d.GetNextMessage()
		if err != nil {
			log.Fatalf("Error getting next message: %s", err)
		}

		response := ""

		switch message.Type {
		case "setup":
			response, err = d.handleSetup(message.Data)
		case "move":
			response, err = d.handleMove(message.Data)
		default:
			log.Fatalf("Unknown message type: %s", message.Type)
		}

		if err != nil {
			log.Fatalf("Error handling message: %+v err: %s", message, err)
		}

		fmt.Println(response)
	}
}

func (d *AIDriver) handleSetup(message []byte) (string, error) {
	setupMessage := tictactoe.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.Symbol = setupMessage.Symbol
	d.state.Order = setupMessage.Order
	d.state.Opponent = setupMessage.Opponent
	return "OK", nil
}

func (d *AIDriver) handleMove(message []byte) (string, error) {
	moveMessage := tictactoe.MessageMove{}
	err := json.Unmarshal(message, &moveMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.Board = tictactoe.GetBoardFromString(moveMessage.Board)
	d.state.NewMoves = moveMessage.NewMoves
	d.state.AllMoves = append(d.state.AllMoves, moveMessage.NewMoves...)

	move := d.ai.GetMove(*d.state)

	// Add new move to our state immediately - we don't get our own moves in NewMoves
	moveLog := tictactoe.MoveLog{
		Move:  move,
		Order: d.state.Order,
	}
	d.state.AllMoves = append(d.state.AllMoves, moveLog)

	moveJSON, err := json.Marshal(&move)
	if err != nil {
		return "", fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return string(moveJSON), nil
}
