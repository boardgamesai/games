package reversi

import "fmt"

const (
	EventTypeMove = "move"
)

type EventMove struct {
	Order int
	Disc
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d/%s plays %s", e.Order, e.Disc, e.Move)
}
