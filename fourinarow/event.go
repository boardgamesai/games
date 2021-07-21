package fourinarow

import (
	"fmt"
	"strings"

	"github.com/boardgamesai/games/game"
)

const (
	EventTypeSetup = "setup"
	EventTypeMove  = "move"
)

type EventSetupPlayer struct {
	ID    game.PlayerID
	Order int
}

type EventSetup struct {
	Players []EventSetupPlayer
}

func (e EventSetup) String() string {
	players := []string{}
	for _, p := range e.Players {
		players = append(players, fmt.Sprintf("%+v", p))
	}
	return "Setup " + strings.Join(players, ", ")
}

type Coords struct {
	Col int
	Row int
}

type EventMove struct {
	ID        game.PlayerID
	WinCoords []Coords `json:",omitempty"`
	Row       int      // The AI only tells us the column but then we compute the row
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("ID %d plays %s", e.ID, e.Move)
}
