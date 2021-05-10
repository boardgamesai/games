package fourinarow

import (
	"fmt"

	"github.com/boardgamesai/games/game"
)

const (
	EventTypeMove = "move"
)

type EventMove struct {
	ID    game.PlayerID
	Order int
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d plays %s", e.ID, e.Move)
}
