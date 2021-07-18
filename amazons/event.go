package amazons

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
	Color SpaceType
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

type EventMove struct {
	ID game.PlayerID
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("ID %d moves %s", e.ID, e.Move)
}
