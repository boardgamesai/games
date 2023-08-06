package driver

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/game/elements/card"
	"github.com/boardgamesai/games/hearts"
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

		d.PrintResponse(response)
	}
}

func (d *Driver) handleSetup(message []byte) ([]byte, error) {
	setupMessage := hearts.MessageSetup{}
	err := json.Unmarshal(message, &setupMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	d.state = d.initState(setupMessage)

	return d.OkJSON(), nil
}

func (d *Driver) handlePass(message []byte) ([]byte, error) {
	passMessage := hearts.MessagePass{}
	err := json.Unmarshal(message, &passMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	err = d.processNewEvents(passMessage.NewEvents)
	if err != nil {
		return []byte{}, err
	}

	d.state.PassDirection = passMessage.Direction

	// Remove these cards from their hand right away
	// No need to sort, because cards will be passed to us and then we'll sort
	move := d.ai.GetPass(*d.state)
	for _, c := range move.Cards {
		d.state.Hand.Remove(c)
	}

	moveJSON, err := json.Marshal(&move)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return moveJSON, nil
}

func (d *Driver) handlePlay(message []byte) ([]byte, error) {
	playMessage := hearts.MessagePlay{}
	err := json.Unmarshal(message, &playMessage)
	if err != nil {
		return []byte{}, fmt.Errorf("JSON decode failed: %s err: %s", message, err)
	}

	err = d.processNewEvents(playMessage.NewEvents)
	if err != nil {
		return []byte{}, err
	}

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
		return []byte{}, fmt.Errorf("JSON encode failed: %+v err: %s", move, err)
	}
	return moveJSON, nil
}

func (d *Driver) processNewEvents(events []game.Event) error {
	for _, e := range events {

		// Note - we don't have to process EventTypePlay - everything we need for state will be set up
		// when we get the move message
		switch e.Type {
		case hearts.EventTypeDeal:
			dealEvent := hearts.EventDeal{}
			err := json.Unmarshal(e.Data, &dealEvent)
			if err != nil {
				return err
			}

			// If we got a hand, we know its ours (we don't get hands for other players)
			d.state.Hand = dealEvent.Hand

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

		case hearts.EventTypeScoreTrick:
			scoreEvent := hearts.EventScoreTrick{}
			err := json.Unmarshal(e.Data, &scoreEvent)
			if err != nil {
				return err
			}

			d.state.CurrentRound[scoreEvent.ID] += scoreEvent.Score

		case hearts.EventTypeScoreRound:
			scoreRoundEvent := hearts.EventScoreRound{}
			err := json.Unmarshal(e.Data, &scoreRoundEvent)
			if err != nil {
				return err
			}

			// Have to translate the player IDs to pointers
			round := map[*hearts.Player]int{}
			for playerID, score := range scoreRoundEvent.RoundScores {
				for _, player := range d.state.Players {
					if player.ID == playerID {
						round[&player] = score
						break
					}
				}
			}

			d.state.Scores.AddRound(round)

			// Reset the current round
			d.state.CurrentRound = map[game.PlayerID]int{}
			d.state.Trick = []card.Card{}
			d.state.TrickCount = 0
			d.state.HeartsBroken = false
		}
	}

	d.state.AddEvents(events)

	return nil
}

func (d *Driver) initState(message hearts.MessageSetup) *State {
	players := []hearts.Player{}
	for _, p := range message.Players {
		players = append(players, *p)
	}

	return &State{
		State: game.State{
			ID: message.ID,
		},
		Position:     message.Position,
		Players:      players,
		Trick:        []card.Card{},
		Scores:       hearts.NewScores(),
		Hand:         hearts.Hand{},
		CurrentRound: map[game.PlayerID]int{},
	}
}
