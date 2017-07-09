package tictactoe

type MessageSetup struct {
	Symbol  string
	Order   int
	Players []*Player
}

type MessageMove struct {
	Board    string
	NewMoves []MoveLog
}
