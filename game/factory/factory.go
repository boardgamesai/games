package factory

import (
	"fmt"

	"github.com/boardgamesai/games/fourinarow"
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts"
	"github.com/boardgamesai/games/reversi"
	"github.com/boardgamesai/games/tictactoe"
)

func New(gameName string) (game.Playable, error) {
	var g game.Playable
	var err error

	switch gameName {
	case "fourinarow":
		g = fourinarow.New()
	case "hearts":
		g = hearts.New()
	case "reversi":
		g = reversi.New()
	case "tictactoe":
		g = tictactoe.New()
	default:
		err = fmt.Errorf("unknown game: %s", gameName)
	}

	return g, err
}
