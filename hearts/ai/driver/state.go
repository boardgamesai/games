package driver

import (
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/game/elements/card"
	"github.com/boardgamesai/games/hearts"
)

type State struct {
	Position     int // Table position
	Players      []hearts.Player
	Trick        []card.Card
	TrickCount   int
	HeartsBroken bool
	CurrentRound map[game.PlayerID]int
	*hearts.Scores
	hearts.Hand
	hearts.PassDirection
	game.State
}
