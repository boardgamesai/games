package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boardgamesai/games/amazons"
	"github.com/boardgamesai/games/game"
)

type AIDriver struct {
	game.AIDriver
	state  *State
	ai     amazonsAI
	colors map[game.PlayerID]amazons.SpaceType
}

func New(ai amazonsAI) *AIDriver {
	return &AIDriver{
		state: &State{
			Board: amazons.NewBoard(),
		},
		ai:     ai,
		colors: map[game.PlayerID]amazons.SpaceType{},
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
	setupMessage := amazons.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.ID = setupMessage.ID
	d.state.Color = setupMessage.Color
	d.state.Order = setupMessage.Order
	d.state.Opponent = setupMessage.Opponent

	d.colors[d.state.ID] = d.state.Color
	d.colors[d.state.Opponent.ID] = d.state.Opponent.Color

	return d.OkJSON(), nil
}

func (d *AIDriver) handleMove(message []byte) ([]byte, error) {
	moveMessage := amazons.MessageMove{}
	err := json.Unmarshal(message, &moveMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	// Apply all our move events to keep the board up to date
	for _, event := range moveMessage.NewEvents {
		if event.Type == amazons.EventTypeMove {
			e := amazons.EventMove{}
			json.Unmarshal(event.Data, &e)
			d.state.Board.ApplyMove(d.colors[e.ID], e.Move)
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
