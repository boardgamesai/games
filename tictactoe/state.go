package tictactoe

type State struct {
	Symbol   string // Symbol of this player
	Order    int    // Order of this player
	Opponent *Player
	Board    *Board
	AllMoves []MoveLog
	NewMoves []MoveLog
}
