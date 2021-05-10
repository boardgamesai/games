package amazons

import (
	"encoding/json"

	"github.com/boardgamesai/games/game"
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

func (c *Comms) Setup(p *Player, other *Player) error {
	message := MessageSetup{
		Color:    p.Color,
		Order:    p.Order,
		ID:       p.ID,
		Opponent: other,
	}
	return p.SendMessageNoResponse(message)
}

func (c *Comms) GetMove(p *Player) (Move, error) {
	move := Move{}

	message := MessageMove{
		NewEvents: c.NewEvents(p.ID),
	}
	responseJSON, err := p.SendMessage(message)
	if err != nil {
		return move, err
	}

	err = json.Unmarshal(responseJSON, &move)
	return move, err
}
