package ulttictactoe

import (
	"fmt"
	"strings"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/tictactoe"
)

const (
	EventTypeSetup = "setup"
	EventTypeMove  = "move"
)

type EventSetupPlayer struct {
	ID     game.PlayerID
	Order  int
	Symbol string
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
	ID          game.PlayerID
	SubFilled   bool             `json:",omitempty"`
	SubWinMoves []Move           `json:",omitempty"`
	WinMoves    []tictactoe.Move `json:",omitempty"`
	Move
}

func (e EventMove) String() string {
	s := fmt.Sprintf("ID %d plays %s", e.ID, e.Move)
	if e.SubFilled {
		s += " (fills subgrid)"
	}
	if len(e.SubWinMoves) > 0 {
		s += " (wins subgrid)"
	}
	return s
}
