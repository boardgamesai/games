package fourinarow

type MessageSetup struct {
	Order   int
	Players []*Player
}

type MessageMove struct {
	Board    string
	NewMoves []MoveLog
}
