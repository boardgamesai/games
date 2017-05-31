package tictactoe

import "fmt"

type MoveLog struct {
	Move
	*Player
}

func (m MoveLog) String() string {
	return fmt.Sprintf("%s plays %s", m.Player, m.Move)
}
