package driver

import "github.com/boardgamesai/games/fourinarow"

type State struct {
	Order    int // Order of this player
	Opponent *fourinarow.Player
	Board    *fourinarow.Board
	AllMoves []fourinarow.MoveLog
	NewMoves []fourinarow.MoveLog
}
