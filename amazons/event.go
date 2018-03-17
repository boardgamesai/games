package amazons

import "fmt"

const (
	EventTypeMove = "move"
)

type EventMove struct {
	Order int
	Color SpaceType
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("P%d/%s moves %s", e.Order, e.Color, e.Move)
}
