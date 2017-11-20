package tictactoe

import "github.com/boardgamesai/games/game"

type MessageSetup struct {
	Symbol   string
	Order    int
	Opponent *Player
}

type MessageMove struct {
	NewEvents []game.Event
}
