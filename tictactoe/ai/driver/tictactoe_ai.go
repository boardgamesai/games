package driver

import "github.com/boardgamesai/games/tictactoe"

type tictactoeAI interface {
	GetMove(state State) tictactoe.Move
}
