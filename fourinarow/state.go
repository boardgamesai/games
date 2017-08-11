package fourinarow

type State struct {
	Order    int // Order of this player
	Players  []Player
	Board    *Board
	AllMoves []MoveLog
	NewMoves []MoveLog
}
