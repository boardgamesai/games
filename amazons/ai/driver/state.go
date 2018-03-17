package driver

import (
	"github.com/boardgamesai/games/amazons"
	"github.com/boardgamesai/games/game"
)

type State struct {
	Color    amazons.SpaceType
	Order    int // Order of this player
	Opponent *amazons.Player
	Board    *amazons.Board
	game.State
}
