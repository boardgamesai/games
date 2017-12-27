package reversi

import (
	"fmt"
	"testing"

	"github.com/boardgamesai/games/game"
)

func getGame(moves map[int][]Move) *Game {
	g := New()
	for i := 1; i <= g.NumPlayers(); i++ {
		g.AddPlayer(fmt.Sprintf("player%d", i), &game.RunnablePlayerMock{})
	}
	g.Comms = &CommsMock{
		moves: moves,
	}
	return g
}

func mv(col, row int) Move {
	return Move{
		Col: col,
		Row: row,
	}
}

func TestSkipTurnGame(t *testing.T) {
	moves := map[int][]Move{
		1: []Move{
			mv(2, 4),
			mv(5, 5),
			mv(2, 2),
			mv(5, 3),
			mv(2, 5),
			mv(3, 5),
			mv(7, 4),
			mv(4, 2),
			mv(6, 3),
			mv(0, 4),
			mv(4, 6),
			mv(7, 6),
			mv(4, 0),
			mv(0, 5),
			mv(5, 7),
			mv(6, 7),
			mv(0, 3),
			mv(6, 2),
			mv(1, 7),
			mv(0, 1),
			mv(3, 0),
			mv(2, 0),
			mv(2, 7),
			mv(4, 7),
			mv(1, 0),
			mv(5, 1),
			mv(1, 6),
			mv(6, 1),
			mv(7, 1),
		},
		2: []Move{
			mv(4, 5),
			mv(2, 3),
			mv(5, 4),
			mv(6, 5),
			mv(6, 4),
			mv(1, 5),
			mv(3, 2),
			mv(5, 2),
			mv(1, 4),
			mv(4, 1),
			mv(7, 5),
			mv(2, 6),
			mv(3, 7),
			mv(5, 6),
			mv(3, 6),
			mv(1, 3),
			mv(7, 2),
			mv(3, 1),
			mv(0, 2),
			mv(1, 2),
			mv(2, 1),
			mv(1, 1),
			mv(0, 7),
			mv(7, 7),
			mv(0, 0),
			mv(7, 3), // just skipped black
			mv(5, 0), // just skipped black
			mv(0, 6),
			mv(6, 6),
			mv(6, 0), // just skipped black
			mv(7, 0),
		},
	}
	g := getGame(moves)

	err := g.Play()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	places := g.Places()
	if places[0].Player.Order != 2 || places[0].Rank != 1 || places[0].Score != 48 {
		t.Errorf("Got incorrect places, first player: %+v", places)
	}

	if places[1].Player.Order != 1 || places[1].Rank != 2 || places[1].Score != 16 {
		t.Errorf("Got incorrect places, second player: %+v", places)
	}
}

func TestGameEndsEarlyWipeout(t *testing.T) {
	moves := map[int][]Move{
		1: []Move{
			mv(5, 3),
			mv(2, 3),
			mv(4, 5),
			mv(6, 3),
			mv(4, 1), // After this white has no discs left
		},
		2: []Move{
			mv(3, 2),
			mv(5, 4),
			mv(5, 2),
			mv(4, 2),
		},
	}
	g := getGame(moves)

	err := g.Play()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	places := g.Places()
	if places[0].Player.Order != 1 || places[0].Rank != 1 || places[0].Score != 13 {
		t.Errorf("Got incorrect places, first player: %+v", places)
	}

	if places[1].Player.Order != 2 || places[1].Rank != 2 || places[1].Score != 0 {
		t.Errorf("Got incorrect places, second player: %+v", places)
	}
}
