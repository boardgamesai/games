package driver

import (
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/tictactoe"
)

type State struct {
	Symbol   string // Symbol of this player
	Order    int    // Order of this player
	Opponent *tictactoe.Player
	Board    *tictactoe.Board
	game.State
}
