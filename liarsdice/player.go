package liarsdice

import (
	"github.com/boardgamesai/games/game"
)

type Player struct {
	game.Player
	Position int // 1-based
}
