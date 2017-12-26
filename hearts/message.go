package hearts

import (
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts/card"
)

type MessageSetup struct {
	Order   int
	Players []*Player
}

type MessagePass struct {
	Direction PassDirection
	NewEvents []game.Event
}

type MessagePlay struct {
	Trick     []card.Card
	NewEvents []game.Event
}
