package tictactoe

type CommsMock struct {
	moves map[int][][]int
	index int
}

func NewCommsMock(moves map[int][][]int) *CommsMock {
	return &CommsMock{
		moves: moves,
		index: -1,
	}
}

func (c *CommsMock) Setup(p *Player, other *Player) error {
	return nil
}

func (c *CommsMock) GetMove(p *Player) (Move, error) {
	if p.Order == 1 {
		c.index++
	}

	move := Move{
		Col: c.moves[p.Order][c.index][0],
		Row: c.moves[p.Order][c.index][1],
	}
	return move, nil
}
