package driver

import (
	"github.com/boardgamesai/games/fourinarow"
	"github.com/boardgamesai/games/game"
)

type State struct {
	Order    int // Order of this player
	Opponent *fourinarow.Player
	Board    *fourinarow.Board
	game.State
}
