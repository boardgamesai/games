package hearts

import (
	"encoding/json"
	"fmt"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts/card"
)

type Player struct {
	game.Player
	Hand
}

func (p *Player) Setup(players []*Player) error {
	message := MessageSetup{
		Order:   p.Order,
		Players: players,
	}
	return p.SendMessageNoResponse(message)
}

func (p *Player) SetHand(newEvents []game.Event) error {
	message := MessageHand{
		Hand:      p.Hand,
		NewEvents: newEvents,
	}
	return p.SendMessageNoResponse(message)
}

func (p *Player) GetPassMove(direction PassDirection, newEvents []game.Event) (PassMove, error) {
	move := PassMove{}

	message := MessagePass{
		Direction: direction,
		NewEvents: newEvents,
	}
	responseJSON, err := p.SendMessage(message)
	if err != nil {
		return move, err
	}

	err = json.Unmarshal(responseJSON, &move)
	return move, err
}

func (p *Player) GetPlayMove(trick []card.Card, newEvents []game.Event) (PlayMove, error) {
	move := PlayMove{}

	message := MessagePlay{
		Trick:     trick,
		NewEvents: newEvents,
	}
	responseJSON, err := p.SendMessage(message)
	if err != nil {
		return move, err
	}

	err = json.Unmarshal(responseJSON, &move)
	return move, err
}

func (p *Player) String() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Order)
}
