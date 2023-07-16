package driver

import "github.com/boardgamesai/games/liarsdice"

type liarsdiceAI interface {
	GetMove(state State) liarsdice.Move
}
