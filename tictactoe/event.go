package tictactoe

import (
	"fmt"

	"github.com/boardgamesai/games/game"
)

const (
	EventTypeMove = "move"
)

type EventMove struct {
	ID       game.PlayerID
	Symbol   string
	WinMoves []Move `json:",omitempty"`
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d plays %s at %s", e.ID, e.Symbol, e.Move)
}
