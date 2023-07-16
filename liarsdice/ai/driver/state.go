package driver

import (
	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/liarsdice"
)

type State struct {
	Position   int // Table position of this player
	Players    []liarsdice.Player
	Dice       []liarsdice.DiceVal
	Bid        liarsdice.DiceVal
	Quantity   int
	Bidder     game.PlayerID
	DiceCounts map[game.PlayerID]int
	DiceShown  map[game.PlayerID][]liarsdice.DiceVal
	game.State
}
