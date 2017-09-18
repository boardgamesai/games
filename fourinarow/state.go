package fourinarow

type State struct {
	Order    int // Order of this player
	Opponent *Player
	Board    *Board
	AllMoves []MoveLog
	NewMoves []MoveLog
}
