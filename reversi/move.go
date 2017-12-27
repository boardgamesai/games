package reversi

import "fmt"

type Move struct {
	Col int
	Row int
}

func (m Move) String() string {
	return fmt.Sprintf("[%d,%d]", m.Col, m.Row)
}
