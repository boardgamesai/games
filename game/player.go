package game

import "fmt"

type Player struct {
	ID    string
	Name  string
	Order int // 1-based
}

func (p *Player) String() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Order)
}
