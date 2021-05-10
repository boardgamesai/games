package fourinarow

import (
	"fmt"
	"testing"

	"github.com/boardgamesai/games/game"
)

func getGame(moves map[int][]int) *Game {
	g := New()
	players := g.Players()
	for i := 0; i < g.MetaData().NumPlayers; i++ {
		players[i].ID = game.PlayerID(i + 1)
		players[i].Name = fmt.Sprintf("player%d", i)
		players[i].Runnable = &game.RunnablePlayerMock{}
	}
	g.Comms = NewCommsMock(moves)
	return g
}

func TestGameWinner(t *testing.T) {
	moves := map[int][]int{
		1: []int{2, 1, 0, 4},
		2: []int{3, 3, 3, 3},
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
	moves := map[int][]int{
		1: []int{1, 0, 1, 0, 1, 0, 2, 3, 2, 3, 2, 3, 5, 4, 5, 4, 5, 4, 6, 6, 6},
		2: []int{0, 1, 0, 1, 0, 1, 3, 2, 3, 2, 3, 2, 4, 5, 4, 5, 4, 5, 6, 6, 6},
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
