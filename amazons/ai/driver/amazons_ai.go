package driver

import "github.com/boardgamesai/games/amazons"

type amazonsAI interface {
	GetMove(state State) amazons.Move
}
