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
	state   *State
	ai      tictactoeAI
	players map[game.PlayerID]string
}

func New(ai tictactoeAI) *AIDriver {
	return &AIDriver{
		state: &State{
			Board: &tictactoe.Board{},
		},
		ai:      ai,
		players: map[game.PlayerID]string{},
	}
}

func (d *AIDriver) Run() {
	d.Setup()

	defer func() {
		if r := recover(); r != nil {
			// This is a panic, trap the error msg and return here
			err := game.DQError{
				ID:   d.state.ID,
				Type: game.DQTypeRuntime,
				Msg:  fmt.Sprintf("%s", r),
			}
			d.PrintErrorResponse(&err)
		}
	}()

	for {
		message, err := d.GetNextMessage()
		if err != nil {
			log.Fatalf("Error getting next message: %s", err)
		}

		var response []byte

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

		d.PrintResponse(response)
	}
}

func (d *AIDriver) handleSetup(message []byte) ([]byte, error) {
	setupMessage := tictactoe.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.Symbol = setupMessage.Symbol
	d.state.ID = setupMessage.ID
	d.state.Order = setupMessage.Order
	d.state.Opponent = setupMessage.Opponent

	d.players[d.state.ID] = d.state.Symbol
	d.players[d.state.Opponent.ID] = d.state.Opponent.Symbol

	return d.OkJSON(), nil
}

func (d *AIDriver) handleMove(message []byte) ([]byte, error) {
	moveMessage := tictactoe.MessageMove{}
	if err := json.Unmarshal(message, &moveMessage); err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	// Apply all our move events to keep the board up to date
	for _, event := range moveMessage.NewEvents {
		if event.Type == tictactoe.EventTypeMove { // This is our only event type, at least for now
			e := tictactoe.EventMove{}
			json.Unmarshal(event.Data, &e)
			d.state.Board.ApplyMove(d.players[e.ID], e.Move)
		}
	}
	d.state.AddEvents(moveMessage.NewEvents)

	move := d.ai.GetMove(*d.state)
	moveJSON, err := json.Marshal(&move)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return moveJSON, nil
}
