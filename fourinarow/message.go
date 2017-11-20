package fourinarow

import "github.com/boardgamesai/games/game"

type MessageSetup struct {
	Order    int
	Opponent *Player
}

type MessageMove struct {
	NewEvents []game.Event
}
