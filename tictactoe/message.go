package tictactoe

const (
	MessageTypeSetup = "setup"
	MessageTypeMove  = "move"
)

type MessageSetup struct {
	Symbol   string
	Order    int
	Opponent *Player
}

type MessageMove struct {
	Board    string
	NewMoves []MoveLog
}
