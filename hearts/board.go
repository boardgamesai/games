package hearts

import "github.com/boardgamesai/games/game/elements/card"

type Board struct {
	Deck   card.StandardDeck
	Hands  map[*Player]*Hand
	Scores *Scores
}

func NewBoard(players []*Player) *Board {
	hands := map[*Player]*Hand{}
	for _, p := range players {
		hands[p] = &Hand{}
	}

	return &Board{
		Deck:   card.NewStandardDeck(),
		Hands:  hands,
		Scores: NewScores(),
	}
}
