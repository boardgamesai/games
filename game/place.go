package game

import "fmt"

type Place struct {
	Player Player
	Rank   int
	Tie    bool
	Score  int
}

func (p Place) String() string {
	return fmt.Sprintf("player: %s rank: %d tie: %t score: %d", p.Player.Name, p.Rank, p.Tie, p.Score)
}
