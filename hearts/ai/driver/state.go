package driver

import (
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts"
	"github.com/boardgamesai/games/hearts/card"
)

type State struct {
	Position     int // Table position
	Players      []hearts.Player
	Trick        []card.Card
	TrickCount   int
	HeartsBroken bool
	RoundScores  map[*hearts.Player]int
	hearts.Hand
	hearts.PassDirection
	game.State
}
