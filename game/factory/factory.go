package factory

import (
	"fmt"

	"github.com/boardgamesai/games/amazons"
	"github.com/boardgamesai/games/fourinarow"
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts"
	"github.com/boardgamesai/games/reversi"
	"github.com/boardgamesai/games/tictactoe"
	"github.com/boardgamesai/games/ulttictactoe"
)

func New(gameName game.Name) (game.Playable, error) {
	var g game.Playable
	var err error

	switch gameName {
	case game.Amazons:
		g = amazons.New()
	case game.FourInARow:
		g = fourinarow.New()
	case game.Hearts:
		g = hearts.New()
	case game.Reversi:
		g = reversi.New()
	case game.TicTacToe:
		g = tictactoe.New()
	case game.UltTicTacToe:
		g = ulttictactoe.New()
	default:
		err = fmt.Errorf("unknown game: %s", gameName)
	}

	return g, err
}
