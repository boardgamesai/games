package driver

import "github.com/boardgamesai/games/reversi"

type reversiAI interface {
	GetMove(state State) reversi.Move
}
