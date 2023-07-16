package liarsdice

import (
	"fmt"

	"github.com/boardgamesai/games/game/elements/dice"
)

type DiceVal int

func (dv DiceVal) String() string {
	if dv == Star {
		return "★"
	}
	return fmt.Sprintf("%d", dv)
}

const Star = DiceVal(1)

var diceVals = []DiceVal{Star, 2, 3, 4, 5, 6}

// Convenience type so we don't have to use generic notation everywhere
type Dice struct {
	*dice.Dice[DiceVal]
}
