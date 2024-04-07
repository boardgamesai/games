package fourinarow

import "github.com/boardgamesai/games/game"

type Player struct {
	game.Player
	Order int // 1-based
}

func NewPlayer() *Player {
	return &Player{
		Player: game.Player{},
	}
}

func (p *Player) BasePlayer() *game.Player {
	return &p.Player
}
