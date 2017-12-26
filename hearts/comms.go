package hearts

import (
	"encoding/json"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts/card"
)

type Comms struct {
	game.Comms
}

func NewComms(g *Game) *Comms {
	return &Comms{
		Comms: game.Comms{
			EventLog: &(g.EventLog),
		},
	}
}

func (c *Comms) Setup(p *Player, players []*Player) error {
	message := MessageSetup{
		Order:   p.Order,
		Players: players,
	}
	return p.SendMessageNoResponse(message)
}

func (c *Comms) GetPassMove(p *Player, direction PassDirection) (PassMove, error) {
	move := PassMove{}

	message := MessagePass{
		Direction: direction,
		NewEvents: c.NewEvents(p.Order),
	}
	responseJSON, err := p.SendMessage(message)
	if err != nil {
		return move, err
	}

	err = json.Unmarshal(responseJSON, &move)
	return move, err
}

func (c *Comms) GetPlayMove(p *Player, trick []card.Card) (PlayMove, error) {
	move := PlayMove{}

	message := MessagePlay{
		Trick:     trick,
		NewEvents: c.NewEvents(p.Order),
	}
	responseJSON, err := p.SendMessage(message)
	if err != nil {
		return move, err
	}

	err = json.Unmarshal(responseJSON, &move)
	return move, err
}
