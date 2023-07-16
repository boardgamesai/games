package liarsdice

import (
	"reflect"
	"testing"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/game/elements/dice"
)

func newPlayer(ID int) *Player {
	return &Player{
		Player: game.Player{
			ID: game.PlayerID(ID),
		},
		Position: ID,
	}
}

func TestIsValidMove(t *testing.T) {
	tests := []struct {
		prevBid      DiceVal
		prevQuantity int
		newBid       DiceVal
		newQuantity  int
		expected     error
	}{
		{0, 0, 1, 1, nil},
		{0, 0, 1, 0, OutOfBoundsError{}},
		{0, 0, 0, 1, OutOfBoundsError{}},
		{0, 0, 1, 20, nil},
		{0, 0, 1, 21, OutOfBoundsError{}},
		{0, 0, 6, 1, nil},
		{0, 0, 7, 1, OutOfBoundsError{}},
		{Star, 2, Star, 1, IllegalMoveError{}},
		{Star, 2, Star, 2, IllegalMoveError{}},
		{Star, 2, Star, 3, nil},
		{6, 5, Star, 2, IllegalMoveError{}},
		{6, 5, Star, 3, nil},
		{6, 5, Star, 9, nil},
		{Star, 2, 2, 3, IllegalMoveError{}},
		{Star, 2, 2, 4, nil},
		{6, 5, 6, 4, IllegalMoveError{}},
		{6, 5, 6, 5, IllegalMoveError{}},
		{6, 5, 6, 6, nil},
		{6, 5, 3, 5, IllegalMoveError{}},
	}

	players := []*Player{
		newPlayer(1),
		newPlayer(2),
		newPlayer(3),
		newPlayer(4),
	}

	for _, test := range tests {
		b := NewBoard(players)
		b.Bid = test.prevBid
		b.Quantity = test.prevQuantity

		m := Move{
			Bid:      test.newBid,
			Quantity: test.newQuantity,
		}
		err := b.IsValidMove(m)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expected) {
			t.Errorf("TestIsValidMove prevBid: %s prevQuantity: %d newBid: %s newQuantity: %d expected: %s got: %s", test.prevBid, test.prevQuantity, test.newBid, test.newQuantity, test.expected, err)
		}
	}
}

func TestIsValidShow(t *testing.T) {
	tests := []struct {
		playerDice []DiceVal
		showDice   []DiceVal
		expected   error
	}{
		{[]DiceVal{2, 3, 4}, []DiceVal{}, nil},
		{[]DiceVal{2, 3, 4}, []DiceVal{2}, nil},
		{[]DiceVal{2, 3, 4}, []DiceVal{5}, IllegalMoveError{}},
		{[]DiceVal{2, 3, Star, 4, Star}, []DiceVal{Star, Star}, nil},
		{[]DiceVal{2, 3, Star, 4, Star}, []DiceVal{Star, Star, 3}, nil},
		{[]DiceVal{2, 3, Star, 4, Star}, []DiceVal{Star, Star, 6}, IllegalMoveError{}},
		{[]DiceVal{5, 5, 5, 5}, []DiceVal{5, 5, 5}, nil},
		{[]DiceVal{5, 5, 5, 5}, []DiceVal{5, 5, 5, 5}, IllegalMoveError{}},
		{[]DiceVal{6, 4, 5, 4, Star}, []DiceVal{4, 4, 3}, IllegalMoveError{}},
		{[]DiceVal{6, 4, 5, 4, Star}, []DiceVal{4, 4}, nil},
	}

	for _, test := range tests {
		player := newPlayer(1)
		b := NewBoard([]*Player{player})

		d := Dice{dice.New(0, diceVals)}
		for _, die := range test.playerDice {
			d.Add(die)
		}
		b.DiceHidden[player] = &d

		m := Move{
			ShowDice: test.showDice,
		}
		err := b.IsValidShow(m, player)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expected) {
			t.Errorf("TestIsValidShow playerDice: %s showDice: %s expected: %s", test.playerDice, test.showDice, test.expected)
		}
	}
}

func TestMoveShownDice(t *testing.T) {
	tests := []struct {
		playerDice    []DiceVal
		showDice      []DiceVal
		newPlayerDice []DiceVal
	}{
		{[]DiceVal{2, 3, 4}, []DiceVal{2}, []DiceVal{3, 4}},
		{[]DiceVal{2, 3, 4}, []DiceVal{3}, []DiceVal{2, 4}},
		{[]DiceVal{2, 3, 4}, []DiceVal{4}, []DiceVal{2, 3}},
		{[]DiceVal{6, 4, 5, 4, Star}, []DiceVal{4, 4}, []DiceVal{6, 5, Star}},
		{[]DiceVal{5, 5, 5, 5}, []DiceVal{5, 5, 5}, []DiceVal{5}},
		{[]DiceVal{4, 4, 6, Star, 3}, []DiceVal{4, 4}, []DiceVal{6, Star, 3}},
	}

	for _, test := range tests {
		player := newPlayer(1)
		b := NewBoard([]*Player{player})

		d := Dice{dice.New(0, diceVals)}
		for _, die := range test.playerDice {
			d.Add(die)
		}
		b.DiceHidden[player] = &d

		m := Move{
			ShowDice: test.showDice,
		}
		b.moveShownDice(m, player)

		if b.DiceHidden[player].Count() != len(test.newPlayerDice) {
			t.Errorf("TestApplyBidMove expected player to have %d dice, had %d", len(test.newPlayerDice), b.DiceHidden[player].Count())
		} else {
			for i := 0; i < len(test.newPlayerDice); i++ {
				if b.DiceHidden[player].Values[i] != test.newPlayerDice[i] {
					t.Errorf("TestApplyBidMove die value mismatch, expected %s saw %s", test.newPlayerDice[i], b.DiceHidden[player].Values[i])
					break
				}
			}
		}
	}
}

func TestDiceChanges(t *testing.T) {
	tests := []struct {
		bid                 int
		actual              int
		expBidderChange     int
		expChallengerChange int
	}{
		{5, 6, 0, -1},
		{5, 7, 0, -2},
		{5, 10, 0, -5},
		{5, 14, 0, -5}, // They only have 5 dice
		{8, 7, -1, 0},
		{8, 6, -2, 0},
		{8, 3, -5, 0},
		{8, 0, -5, 0}, // They only have 5 dice
		{1, 0, -1, 0},
		{4, 4, 1, -1},
	}

	for _, test := range tests {
		bidder := newPlayer(1)
		challenger := newPlayer(2)
		b := NewBoard([]*Player{bidder, challenger})
		b.Quantity = test.bid
		b.Bidder = bidder

		changes := b.calculateDiceChanges(test.actual, challenger)
		if changes[bidder] != test.expBidderChange || changes[challenger] != test.expChallengerChange {
			t.Errorf("TestDiceChanges got %d, %d expected %d, %d", changes[bidder], changes[challenger], test.expBidderChange, test.expChallengerChange)
			continue
		}

		b.applyDiceChanges(changes)
		if b.DiceHidden[bidder].Count() != 5+test.expBidderChange {
			t.Errorf("TestDiceChanges expected bidder to have %d dice, had %d", test.expBidderChange, b.DiceHidden[bidder].Count())
		}
		if b.DiceHidden[challenger].Count() != 5+test.expChallengerChange {
			t.Errorf("TestDiceChanges expected challenger to have %d dice, had %d", test.expChallengerChange, b.DiceHidden[challenger].Count())
		}
	}
}
