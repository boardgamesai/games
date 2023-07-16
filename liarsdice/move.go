package liarsdice

import (
	"fmt"
)

type Move struct {
	Challenge bool      `json:",omitempty"`
	Bid       DiceVal   `json:",omitempty"`
	Quantity  int       `json:",omitempty"`
	ShowDice  []DiceVal `json:",omitempty"`
}

func (m Move) String() string {
	s := ""
	if m.Challenge {
		s = "challenge"
	} else {
		s = fmt.Sprintf("bid %d %ss", m.Quantity, m.Bid)
		if len(m.ShowDice) > 0 {
			s += fmt.Sprintf(", show %s", m.ShowDice)
		}
	}

	return s
}
