package tictactoe

type Bot interface {
	GetMove(symbol string, board Board) Move
}
