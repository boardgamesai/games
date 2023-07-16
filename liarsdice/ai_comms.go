package liarsdice

type AIComms interface {
	Setup(p *Player, players []*Player) error
	GetMove(p *Player) (Move, error)
}
