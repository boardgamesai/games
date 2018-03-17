package amazons

type CommsMock struct {
	moves map[int][]Move
}

func (c *CommsMock) Setup(p *Player, other *Player) error {
	return nil
}

func (c *CommsMock) GetMove(p *Player) (Move, error) {
	m := c.moves[p.Order][0]
	c.moves[p.Order] = c.moves[p.Order][1:]
	return m, nil
}
