package tictactoe

import (
	"encoding/json"
	"fmt"

	"github.com/boardgamesai/games/game"
)

type Player struct {
	game.Player
	Symbol string // "X" or "O"
}

func (p *Player) Setup(other *Player) error {
	message := MessageSetup{
		Symbol:   p.Symbol,
		Order:    p.Order,
		Opponent: other,
	}
	response, err := p.SendMessage(message)
	if err != nil {
		return err
	}
	if string(response) != "OK" {
		return fmt.Errorf("Got non-OK response when setting up player: %s err: %s", p.Name, err)
	}

	return nil
}

func (p *Player) GetMove(newEvents []game.Event) (Move, error) {
	move := Move{}

	message := MessageMove{
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
	return fmt.Sprintf("%s (%s)", p.Name, p.Symbol)
}
