package driver

import "github.com/boardgamesai/games/hearts"

type heartsAI interface {
	GetPass(state State) hearts.PassMove
	GetPlay(state State) hearts.PlayMove
}
