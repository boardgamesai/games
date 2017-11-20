package main

import (
	"github.com/boardgamesai/games/hearts"
	"github.com/boardgamesai/games/hearts/ai/driver"
)

type AI struct{}

func (ai *AI) GetPass(state driver.State) hearts.PassMove {
	// TODO
	return hearts.PassMove{}
}

func (ai *AI) GetPlay(state driver.State) hearts.PlayMove {
	// TODO
	return hearts.PlayMove{}
}
