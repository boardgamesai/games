package fourinarow

import "fmt"

const (
	EventTypeMove = "move"
)

type EventMove struct {
	Order int
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d plays %s", e.Order, e.Move)
}
