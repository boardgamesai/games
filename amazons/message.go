package amazons

import "github.com/boardgamesai/games/game"

type MessageSetup struct {
	Color    SpaceType
	Order    int
	ID       game.PlayerID
	Opponent *Player
}

type MessageMove struct {
	NewEvents []game.Event
}
