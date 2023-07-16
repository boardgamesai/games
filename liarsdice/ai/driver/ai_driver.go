package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/liarsdice"
)

type Driver struct {
	game.AIDriver
	state *State
	ai    liarsdiceAI
}

func New(ai liarsdiceAI) *Driver {
	return &Driver{
		ai:    ai,
		state: &State{},
	}
}

func (d *Driver) Run() {
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

func (d *Driver) handleSetup(message []byte) ([]byte, error) {
	setupMessage := liarsdice.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state = d.initState(setupMessage)

	return d.OkJSON(), nil
}

func (d *Driver) handleMove(message []byte) ([]byte, error) {
	moveMessage := liarsdice.MessageMove{}
	err := json.Unmarshal(message, &moveMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	// Update our state based on all the new events
	for _, e := range moveMessage.NewEvents {
		switch e.Type {
		case liarsdice.EventTypeRoll:
			// We will only get a roll if it's ours, no need to check IDs
			eventRoll := liarsdice.EventRoll{}
			if err := json.Unmarshal(e.Data, &eventRoll); err != nil {
				return []byte{}, err
			}

			d.state.Dice = eventRoll.Dice

		case liarsdice.EventTypeMove:
			// This is a non-challege move, challenges get their own event
			eventMove := liarsdice.EventMove{}
			if err := json.Unmarshal(e.Data, &eventMove); err != nil {
				return []byte{}, err
			}

			d.state.Bid = eventMove.Bid
			d.state.Quantity = eventMove.Quantity
			d.state.Bidder = eventMove.ID

			if len(eventMove.ShowDice) > 0 {
				d.state.DiceShown[eventMove.ID] = append(d.state.DiceShown[eventMove.ID], eventMove.ShowDice...)
			}

		case liarsdice.EventTypeChallenge:
			eventChallenge := liarsdice.EventChallenge{}
			if err := json.Unmarshal(e.Data, &eventChallenge); err != nil {
				return []byte{}, err
			}

			for ID, change := range eventChallenge.DiceChange {
				d.state.DiceCounts[ID] += change
			}

			// Reset the round
			d.state.Bid = 0
			d.state.Quantity = 0
			d.state.Bidder = 0
			for _, p := range d.state.Players {
				d.state.DiceShown[p.ID] = []liarsdice.DiceVal{}
			}
		}
	}
	d.state.AddEvents(moveMessage.NewEvents)

	move := d.ai.GetMove(*(d.state))
	moveJSON, err := json.Marshal(&move)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return moveJSON, nil
}

func (d *Driver) initState(message liarsdice.MessageSetup) *State {
	players := []liarsdice.Player{}
	diceCounts := map[game.PlayerID]int{}
	diceShown := map[game.PlayerID][]liarsdice.DiceVal{}

	for _, p := range message.Players {
		players = append(players, *p)
		diceCounts[p.ID] = 5
		diceShown[p.ID] = []liarsdice.DiceVal{}
	}

	return &State{
		State: game.State{
			ID: message.ID,
		},
		Position:   message.Position,
		Players:    players,
		Dice:       []liarsdice.DiceVal{},
		DiceCounts: diceCounts,
		DiceShown:  diceShown,
	}
}
