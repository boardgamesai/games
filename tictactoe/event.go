package tictactoe

import "fmt"

const (
	EventTypeMove = "move"
)

type EventMove struct {
	Order  int
	Symbol string
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d plays %s at %s", e.Order, e.Symbol, e.Move)
}
