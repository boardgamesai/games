package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts"
	"github.com/boardgamesai/games/hearts/card"
)

type Driver struct {
	game.AIDriver
	state *State
	ai    heartsAI
}

func New(ai heartsAI) *Driver {
	return &Driver{
		state: &State{},
		ai:    ai,
	}
}

func (d *Driver) Run() {
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
		case "hand":
			response, err = d.handleHand(message.Data)
		case "pass":
			response, err = d.handlePass(message.Data)
		case "play":
			response, err = d.handlePlay(message.Data)
		default:
			log.Fatalf("Unknown message type: %s", message.Type)
		}

		if err != nil {
			log.Fatalf("Error handling message: %+v err: %s", message, err)
		}

		fmt.Println(response)
	}
}

func (d *Driver) handleSetup(message []byte) (string, error) {
	setupMessage := hearts.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.Order = setupMessage.Order
	d.state.Players = []hearts.Player{}
	return "OK", nil
}

func (d *Driver) handleHand(message []byte) (string, error) {
	handMessage := hearts.MessageHand{}
	err := json.Unmarshal(message, &handMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.AddEvents(handMessage.NewEvents)
	d.state.Hand = handMessage.Hand
	d.state.Trick = []card.Card{}
	d.state.TrickCount = 0
	d.state.HeartsBroken = false

	return "OK", nil
}

func (d *Driver) handlePass(message []byte) (string, error) {
	passMessage := hearts.MessagePass{}
	err := json.Unmarshal(message, &passMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.AddEvents(passMessage.NewEvents)
	d.state.PassDirection = passMessage.Direction

	move := d.ai.GetPass(*d.state)
	moveJSON, err := json.Marshal(&move)
	if err != nil {
		return "", fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return string(moveJSON), nil
}

func (d *Driver) handlePlay(message []byte) (string, error) {
	playMessage := hearts.MessagePlay{}
	err := json.Unmarshal(message, &playMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state.AddEvents(playMessage.NewEvents)
	d.state.Trick = playMessage.Trick
	for _, c := range d.state.Trick {
		if c.Suit == card.Hearts {
			d.state.HeartsBroken = true
		}
	}

	move := d.ai.GetPlay(*d.state)

	if !d.state.HeartsBroken && move.Card.Suit == card.Hearts {
		d.state.HeartsBroken = true
	}
	d.state.Hand.Remove(move.Card)
	d.state.TrickCount++

	moveJSON, err := json.Marshal(&move)
	if err != nil {
		return "", fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return string(moveJSON), nil
}
