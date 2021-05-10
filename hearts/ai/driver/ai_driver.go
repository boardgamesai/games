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

	d.state.ID = setupMessage.ID
	d.state.Position = setupMessage.Position
	d.state.Players = []hearts.Player{}
	return "OK", nil
}

func (d *Driver) handlePass(message []byte) (string, error) {
	passMessage := hearts.MessagePass{}
	err := json.Unmarshal(message, &passMessage)
	if err != nil {
		return "", fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	err = d.checkForDealPass(passMessage.NewEvents)
	if err != nil {
		return "", err
	}

	d.state.AddEvents(passMessage.NewEvents)
	d.state.PassDirection = passMessage.Direction

	// Remove these cards from their hand right away
	// No need to sort, because cards will be passed to us and then we'll sort
	move := d.ai.GetPass(*d.state)
	for _, c := range move.Cards {
		d.state.Hand.Remove(c)
	}

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

	err = d.checkForDealPass(playMessage.NewEvents)
	if err != nil {
		return "", err
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

func (d *Driver) checkForDealPass(events []game.Event) error {
	for _, e := range events {

		switch e.Type {
		case hearts.EventTypeDeal:
			dealEvent := hearts.EventDeal{}
			err := json.Unmarshal(e.Data, &dealEvent)
			if err != nil {
				return err
			}

			// If we got a hand, we know its ours (we don't get hands for other players) and we reset
			d.state.Hand = dealEvent.Hand
			d.state.Trick = []card.Card{}
			d.state.TrickCount = 0
			d.state.HeartsBroken = false
		case hearts.EventTypePass:
			passEvent := hearts.EventPass{}
			err := json.Unmarshal(e.Data, &passEvent)
			if err != nil {
				return err
			}

			if passEvent.ToID != d.state.ID {
				continue
			}

			for _, c := range passEvent.Cards {
				d.state.Hand.Add(c)
			}
			d.state.Hand.Sort()
		}
	}

	return nil
}
