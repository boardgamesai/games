package tictactoe

type AIComms interface {
	Setup(p *Player, other *Player) error
	GetMove(p *Player) (Move, error)
}
