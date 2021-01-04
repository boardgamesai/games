package game

import "fmt"

type Player struct {
	Runnable `json:"-"`
	ID       string
	Name     string `json:",omitempty"`
	Order    int    // 1-based
}

func (p *Player) String() string {
	return fmt.Sprintf("%s (%d)", p.Name, p.Order)
}
