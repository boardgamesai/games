package driver

import (
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/ulttictactoe"
)

type State struct {
	Symbol   string // Symbol of this player
	Order    int    // Order of this player
	Opponent *ulttictactoe.Player
	Board    *ulttictactoe.Board
	game.State
}
