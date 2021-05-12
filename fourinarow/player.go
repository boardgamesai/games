package fourinarow

import "github.com/boardgamesai/games/game"

type Player struct {
	game.Player
	Order int // 1-based
}
