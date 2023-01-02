package game

import "fmt"

type PlayerID uint64

type Player struct {
	Runnable `json:"-"`
	ID       PlayerID
	Name     string `json:",omitempty"`
}

func (p *Player) String() string {
	if p.Name == "" {
		return fmt.Sprintf("%d", p.ID)
	}

	return fmt.Sprintf("%s (%d)", p.Name, p.ID)
}
