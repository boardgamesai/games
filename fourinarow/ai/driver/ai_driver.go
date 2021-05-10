package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boardgamesai/games/fourinarow"
	"github.com/boardgamesai/games/game"
)

type AIDriver struct {
	game.AIDriver
	state *State
	ai    fourinarowAI
}

func New(ai fourinarowAI) *AIDriver {
	return &AIDriver{
		state: &State{
			Board: &fourinarow.Board{},
		},
		ai: ai,
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
	setupMessage := fourinarow.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.ID = setupMessage.ID
	d.state.Order = setupMessage.Order
	d.state.Opponent = setupMessage.Opponent
	return "OK", nil
}

func (d *AIDriver) handleMove(message []byte) (string, error) {
	moveMessage := fourinarow.MessageMove{}
	err := json.Unmarshal(message, &moveMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	// Apply all our move events to keep the board up to date
	for _, event := range moveMessage.NewEvents {
		if event.Type == fourinarow.EventTypeMove { // This is our only event type, at least for now
			e := fourinarow.EventMove{}
			json.Unmarshal(event.Data, &e)
			d.state.Board.ApplyMove(e.Order, e.Move)
		}
	}
	d.state.AddEvents(moveMessage.NewEvents)

	move := d.ai.GetMove(*d.state)
	moveJSON, err := json.Marshal(&move)
	if err != nil {
		return "", fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return string(moveJSON), nil
}
