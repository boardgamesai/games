package liarsdice

import (
	"fmt"
	"strings"

	"github.com/boardgamesai/games/game"
)

const (
	EventTypeSetup     = "setup"
	EventTypeMove      = "move"
	EventTypeChallenge = "challenge"
	EventTypeRoll      = "roll"
)

type EventSetupPlayer struct {
	ID       game.PlayerID
	Position int
}

type EventSetup struct {
	Players []EventSetupPlayer
}

func (e EventSetup) String() string {
	players := []string{}
	for _, p := range e.Players {
		players = append(players, fmt.Sprintf("%+v", p))
	}
	return "Setup " + strings.Join(players, ", ")
}

type EventRoll struct {
	ID   game.PlayerID
	Dice []DiceVal
}

func (e EventRoll) String() string {
	return fmt.Sprintf("ID %d rolls %s", e.ID, e.Dice)
}

type EventMove struct {
	ID game.PlayerID
	Move
}

func (e EventMove) String() string {
	return fmt.Sprintf("ID %d %s", e.ID, e.Move)
}

type EventChallenge struct {
	ID             game.PlayerID
	Bid            DiceVal
	ActualQuantity int
	DiceChange     map[game.PlayerID]int
	Eliminated     game.PlayerID `json:",omitempty"`
}

func (e EventChallenge) String() string {
	s := ""
	for p, count := range e.DiceChange {
		s += fmt.Sprintf("%d:%d ", p, count)
	}
	s = strings.TrimSpace(s)
	return fmt.Sprintf("ID %d challenge actual: %d dice change: %s", e.ID, e.ActualQuantity, s)
}
