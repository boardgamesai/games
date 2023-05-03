package hearts

import "github.com/boardgamesai/games/game/elements/card"

type AIComms interface {
	Setup(p *Player, players []*Player) error
	GetPassMove(p *Player, direction PassDirection) (PassMove, error)
	GetPlayMove(p *Player, trick []card.Card) (PlayMove, error)
}
