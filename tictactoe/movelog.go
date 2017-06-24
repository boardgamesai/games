package tictactoe

import "fmt"

type MoveLog struct {
	Move
	Order int
}

func (m MoveLog) String() string {
	return fmt.Sprintf("Player %d plays %s", m.Order, m.Move)
}
