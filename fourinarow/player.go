package fourinarow

import (
	"encoding/json"
	"fmt"

	"github.com/boardgamesai/games/game"
)

type Player struct {
	ID    string
	Name  string
	Order int // 1-based
	game.RunnablePlayer
}

func (p *Player) Setup(g *Game) error {
	// Note that we don't wait for a response here.
	// The other end reads line-by-line and will do the right thing.
	message := MessageSetup{
		Order:   p.Order,
		Players: g.Players,
	}
	messageJSON, err := json.Marshal(&message)
	if err != nil {
		return err
	}

	response, err := p.SendMessage(string(messageJSON))
	if err != nil {
		return err
	}
	if response != "OK" {
		return fmt.Errorf("Got non-OK response when setting up player: %s err: %s", p.Name, err)
	}

	return nil
}

func (p *Player) GetMove(g *Game) (Move, error) {
	message := MessageMove{
		Board:    GetStringFromBoard(g.Board),
		NewMoves: g.GetNewMovesForPlayer(p),
	}
	messageJSON, err := json.Marshal(&message)
	if err != nil {
		return Move{}, err
	}

	responseJSON, err := p.SendMessage(string(messageJSON))
	if err != nil {
		return Move{}, err
	}

	move := Move{}
	err = json.Unmarshal([]byte(responseJSON), &move)
	return move, err
}

func (p *Player) String() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Order)
}
