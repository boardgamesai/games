package fourinarow

type MessageSetup struct {
	Order    int
	Opponent *Player
}

type MessageMove struct {
	Board    string
	NewMoves []MoveLog
}
