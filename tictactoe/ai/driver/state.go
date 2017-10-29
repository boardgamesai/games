package driver

import "github.com/boardgamesai/games/tictactoe"

type State struct {
	Symbol   string // Symbol of this player
	Order    int    // Order of this player
	Opponent *tictactoe.Player
	Board    *tictactoe.Board
	AllMoves []tictactoe.MoveLog
	NewMoves []tictactoe.MoveLog
}
