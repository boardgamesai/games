package amazons

import "github.com/boardgamesai/games/game"

type MessageSetup struct {
	Color    SpaceType
	Order    int
	Opponent *Player
}

type MessageMove struct {
	NewEvents []game.Event
}
