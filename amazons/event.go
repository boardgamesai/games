package amazons

import (
	"fmt"

	"github.com/boardgamesai/games/game"
)

const (
	EventTypeMove = "move"
)

type EventMove struct {
	ID    game.PlayerID
	Color SpaceType
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d/%s moves %s", e.ID, e.Color, e.Move)
}
