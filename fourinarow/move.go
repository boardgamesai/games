package fourinarow

import "fmt"

type Move struct {
	Col int
}

func (m Move) String() string {
	return fmt.Sprintf("%d", m.Col)
}
