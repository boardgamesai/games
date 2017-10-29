package driver

import "github.com/boardgamesai/games/fourinarow"

type fourinarowAI interface {
	GetMove(state State) fourinarow.Move
}
