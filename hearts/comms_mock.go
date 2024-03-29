package hearts

import "github.com/boardgamesai/games/game/elements/card"

type CommsMock struct {
	hands map[int][]string // Assume that hands are set up in the order we want them played
	index int
	count int
}

func (c *CommsMock) Setup(p *Player, players []*Player) error {
	return nil
}

func (c *CommsMock) GetPassMove(p *Player, direction PassDirection) (PassMove, error) {
	hand := c.hands[p.Position]
	move := PassMove{
		Cards: []card.Card{
			card.FromString(hand[0]),
			card.FromString(hand[1]),
			card.FromString(hand[2]),
		},
	}

	return move, nil
}

func (c *CommsMock) GetPlayMove(p *Player, trick []card.Card) (PlayMove, error) {
	move := PlayMove{
		Card: card.FromString(c.hands[p.Position][c.index]),
	}

	c.count++
	if (c.count % 4) == 0 {
		c.index++
	}

	return move, nil
}
