package hearts

import "github.com/boardgamesai/games/game"

type Player struct {
	game.Runnable `json:"-"`
	game.Player
	Hand
}
