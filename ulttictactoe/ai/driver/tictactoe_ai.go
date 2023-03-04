package driver

import "github.com/boardgamesai/games/ulttictactoe"

type ulttictactoeAI interface {
	GetMove(state State) ulttictactoe.Move
}
