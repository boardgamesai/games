package fourinarow

import (
	"fmt"

	"github.com/boardgamesai/games/game"
)

const (
	EventTypeMove = "move"
)

type Coords struct {
	Col int
	Row int
}

type EventMove struct {
	ID        game.PlayerID
	Order     int
	WinCoords []Coords `json:",omitempty"`
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d plays %s", e.ID, e.Move)
}
