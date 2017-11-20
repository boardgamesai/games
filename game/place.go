package game

import "fmt"

type Place struct {
	Player Player
	Rank   int
	Tie    bool
}

func (p Place) String() string {
	return fmt.Sprintf("player: %s rank: %d tie: %t", p.Player.Name, p.Rank, p.Tie)
}
