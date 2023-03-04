package ulttictactoe

import "fmt"

type Move struct {
	Col    int
	Row    int
	SubCol int
	SubRow int
}

func (m Move) String() string {
	return fmt.Sprintf("grid [%d,%d] subgrid [%d,%d]", m.Col, m.Row, m.SubCol, m.SubRow)
}
