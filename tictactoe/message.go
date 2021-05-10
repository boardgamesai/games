package tictactoe

import "github.com/boardgamesai/games/game"

type MessageSetup struct {
	Symbol   string
	Order    int
	ID       game.PlayerID
	Opponent *Player
}

type MessageMove struct {
	NewEvents []game.Event
}
