package liarsdice

import (
	"github.com/boardgamesai/games/game"
)

type MessageSetup struct {
	ID       game.PlayerID
	Position int
	Players  []*Player
}

type MessageMove struct {
	NewEvents []game.Event
}
