package tictactoe

import (
	"fmt"
	"testing"

	"github.com/boardgamesai/games/game"
)

func getGame(moves map[int][][]int) *Game {
	g := New()
	for i := 1; i <= g.NumPlayers(); i++ {
		g.AddPlayer(fmt.Sprintf("player%d", i), &game.RunnablePlayerMock{})
	}
	g.Comms = NewCommsMock(moves)
	return g
}

func TestGameWinner(t *testing.T) {
	moves := map[int][][]int{
		1: [][]int{[]int{1, 2}, []int{2, 2}, []int{2, 1}},
		2: [][]int{[]int{1, 1}, []int{0, 2}, []int{2, 0}},
	}
	g := getGame(moves)

	err := g.Play()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	places := g.Places()
	if places[0].Player.Order != 2 || places[0].Rank != 1 || places[0].Tie {
		t.Errorf("Got incorrect places, player 1: %+v", places)
	}
	if places[1].Player.Order != 1 || places[1].Rank != 2 || places[1].Tie {
		t.Errorf("Got incorrect places, player 2: %+v", places)
	}
}

func TestGameTie(t *testing.T) {
	moves := map[int][][]int{
		1: [][]int{[]int{1, 1}, []int{2, 1}, []int{0, 2}, []int{1, 2}, []int{0, 0}},
		2: [][]int{[]int{2, 2}, []int{0, 1}, []int{2, 0}, []int{1, 0}},
	}
	g := getGame(moves)

	err := g.Play()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	places := g.Places()
	if places[0].Rank != 1 || !places[0].Tie {
		t.Errorf("Got incorrect places, player 1: %+v", places)
	}
	if places[1].Rank != 1 || !places[1].Tie {
		t.Errorf("Got incorrect places, player 2: %+v", places)
	}
}