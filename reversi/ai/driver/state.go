package driver

import (
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/reversi"
)

type State struct {
	reversi.Disc
	Order    int // Order of this player
	Opponent *reversi.Player
	Board    *reversi.Board
	game.State
}
