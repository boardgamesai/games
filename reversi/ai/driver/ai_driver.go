package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/reversi"
)

type AIDriver struct {
	game.AIDriver
	state *State
	ai    reversiAI
	discs map[game.PlayerID]reversi.Disc
}

func New(ai reversiAI) *AIDriver {
	return &AIDriver{
		state: &State{
			Board: reversi.NewBoard(),
		},
		ai:    ai,
		discs: map[game.PlayerID]reversi.Disc{},
	}
}

func (d *AIDriver) Run() {
	d.Setup()
	defer d.HandlePanic(d.state.ID)

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
	setupMessage := reversi.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.Disc = setupMessage.Disc
	d.state.ID = setupMessage.ID
	d.state.Order = setupMessage.Order
	d.state.Opponent = setupMessage.Opponent

	d.discs[d.state.ID] = d.state.Disc
	d.discs[d.state.Opponent.ID] = d.state.Opponent.Disc

	return d.OkJSON(), nil
}

func (d *AIDriver) handleMove(message []byte) ([]byte, error) {
	moveMessage := reversi.MessageMove{}
	err := json.Unmarshal(message, &moveMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	// Apply all our move events to keep the board up to date
	for _, event := range moveMessage.NewEvents {
		if event.Type == reversi.EventTypeMove {
			e := reversi.EventMove{}
			json.Unmarshal(event.Data, &e)
			d.state.Board.ApplyMove(d.discs[e.ID], e.Move)
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
