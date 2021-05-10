package reversi

import (
	"fmt"

	"github.com/boardgamesai/games/game"
)

const (
	EventTypeMove = "move"
)

type EventMove struct {
	ID game.PlayerID
	Disc
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d/%s plays %s", e.ID, e.Disc, e.Move)
}
