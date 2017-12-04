package tictactoe

import "github.com/boardgamesai/games/game"

type Player struct {
	game.Runnable `json:"-"`
	game.Player
	Symbol string // "X" or "O"
}
