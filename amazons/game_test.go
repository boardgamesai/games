package amazons

import (
	"fmt"
	"testing"

	"github.com/boardgamesai/games/game"
)

func getGame(moves map[int][]Move) *Game {
	g := New()
	players := g.Players()
	for i := 0; i < g.MetaData().NumPlayers; i++ {
		players[i].ID = game.PlayerID(i + 1)
		players[i].Name = fmt.Sprintf("player%d", i)
		players[i].Runnable = &game.RunnablePlayerMock{}
	}
	g.Comms = &CommsMock{
		moves: moves,
	}
	return g
}

func mv(from, to, arrow Space) Move {
	return Move{
		From:  from,
		To:    to,
		Arrow: arrow,
	}
}

func sp(col, row int) Space {
	return Space{
		Col: col,
		Row: row,
	}
}

func TestShortestPossibleGame(t *testing.T) {
	moves := map[int][]Move{
		1: []Move{
			mv(sp(0, 3), sp(0, 0), sp(1, 1)),
			mv(sp(3, 0), sp(1, 0), sp(2, 1)),
			mv(sp(6, 0), sp(2, 0), sp(3, 1)),
			mv(sp(9, 3), sp(6, 0), sp(9, 0)),
			mv(sp(6, 0), sp(3, 0), sp(4, 1)),
		},
		2: []Move{
			mv(sp(0, 6), sp(0, 1), sp(0, 9)),
			mv(sp(9, 6), sp(9, 5), sp(9, 9)),
			mv(sp(9, 5), sp(9, 6), sp(9, 8)),
			mv(sp(9, 6), sp(9, 5), sp(9, 7)),
			mv(sp(9, 5), sp(8, 4), sp(4, 0)),
		},
	}
	expectedMoves := len(moves[1]) + len(moves[2])

	g := getGame(moves)
	err := g.Play()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if len(g.Events()) != expectedMoves {
		t.Errorf("Expected %d moves, got %d", expectedMoves, len(g.Events()))
	}

	places := g.Places()
	if places[0].Player.ID != g.players[1].ID || places[0].Rank != 1 {
		t.Errorf("Got incorrect places, first player: %+v", places[0])
	}

	if places[1].Player.ID != g.players[0].ID || places[1].Rank != 2 {
		t.Errorf("Got incorrect places, second player: %+v", places[1])
	}
}

func TestPlayGame(t *testing.T) {
	moves := map[int][]Move{
		1: []Move{
			mv(sp(3, 0), sp(3, 3), sp(3, 8)),
			mv(sp(9, 3), sp(4, 8), sp(4, 9)),
			mv(sp(4, 8), sp(3, 7), sp(2, 6)),
			mv(sp(3, 3), sp(2, 3), sp(2, 4)),
			mv(sp(3, 7), sp(5, 7), sp(3, 7)),
			mv(sp(6, 0), sp(4, 2), sp(4, 8)),
			mv(sp(0, 3), sp(0, 1), sp(0, 3)),
			mv(sp(2, 3), sp(4, 5), sp(2, 3)),
			mv(sp(4, 2), sp(8, 2), sp(6, 2)),
			mv(sp(5, 7), sp(6, 6), sp(6, 4)),
			mv(sp(8, 2), sp(8, 3), sp(5, 6)),
			mv(sp(6, 6), sp(7, 6), sp(7, 2)),
			mv(sp(0, 1), sp(1, 0), sp(3, 2)),
			mv(sp(4, 5), sp(4, 7), sp(4, 6)),
			mv(sp(1, 0), sp(4, 0), sp(4, 2)),
			mv(sp(4, 0), sp(3, 1), sp(2, 1)),
			mv(sp(4, 7), sp(6, 9), sp(8, 7)),
			mv(sp(7, 6), sp(8, 6), sp(9, 6)),
			mv(sp(6, 9), sp(7, 8), sp(7, 7)),
			mv(sp(8, 6), sp(9, 7), sp(9, 8)),
			mv(sp(3, 1), sp(4, 1), sp(3, 1)),
			mv(sp(7, 8), sp(6, 8), sp(7, 8)),
			mv(sp(6, 8), sp(6, 9), sp(5, 9)),
			mv(sp(6, 9), sp(6, 8), sp(6, 9)),
			mv(sp(6, 8), sp(5, 7), sp(6, 8)),
			mv(sp(5, 7), sp(4, 7), sp(5, 7)),
			mv(sp(4, 1), sp(4, 0), sp(3, 0)),
			mv(sp(4, 0), sp(4, 1), sp(4, 0)),
			mv(sp(4, 1), sp(5, 2), sp(4, 1)),
			mv(sp(9, 7), sp(8, 6), sp(9, 7)),
			mv(sp(8, 6), sp(7, 6), sp(8, 6)),
			mv(sp(5, 2), sp(6, 1), sp(5, 2)),
			mv(sp(6, 1), sp(6, 0), sp(6, 1)),
			mv(sp(6, 0), sp(7, 0), sp(6, 0)),
			mv(sp(7, 0), sp(7, 1), sp(7, 0)),
			mv(sp(7, 1), sp(8, 0), sp(7, 1)),
			mv(sp(8, 0), sp(9, 0), sp(8, 0)),
			mv(sp(9, 0), sp(9, 1), sp(9, 0)),
		},
		2: []Move{
			mv(sp(6, 9), sp(6, 1), sp(2, 5)),
			mv(sp(3, 9), sp(1, 7), sp(1, 2)),
			mv(sp(1, 7), sp(1, 4), sp(1, 3)),
			mv(sp(9, 6), sp(4, 6), sp(3, 6)),
			mv(sp(4, 6), sp(5, 6), sp(5, 0)),
			mv(sp(6, 1), sp(6, 3), sp(4, 3)),
			mv(sp(6, 3), sp(6, 2), sp(5, 1)),
			mv(sp(6, 2), sp(6, 3), sp(5, 3)),
			mv(sp(5, 6), sp(5, 5), sp(1, 1)),
			mv(sp(6, 3), sp(8, 5), sp(6, 7)),
			mv(sp(5, 5), sp(5, 4), sp(5, 5)),
			mv(sp(5, 4), sp(3, 4), sp(5, 4)),
			mv(sp(3, 4), sp(3, 5), sp(4, 4)),
			mv(sp(3, 5), sp(3, 3), sp(3, 5)),
			mv(sp(3, 3), sp(2, 2), sp(2, 0)),
			mv(sp(8, 5), sp(8, 6), sp(8, 4)),
			mv(sp(8, 6), sp(9, 7), sp(7, 9)),
			mv(sp(9, 7), sp(8, 8), sp(5, 8)),
			mv(sp(8, 8), sp(8, 9), sp(8, 8)),
			mv(sp(1, 4), sp(0, 4), sp(1, 4)),
			mv(sp(2, 2), sp(3, 3), sp(2, 2)),
			mv(sp(8, 9), sp(9, 9), sp(8, 9)),
			mv(sp(3, 3), sp(3, 4), sp(4, 5)),
			mv(sp(3, 4), sp(3, 3), sp(3, 4)),
			mv(sp(0, 4), sp(0, 5), sp(0, 4)),
			mv(sp(0, 5), sp(1, 5), sp(0, 5)),
			mv(sp(1, 5), sp(1, 6), sp(1, 5)),
			mv(sp(1, 6), sp(0, 7), sp(1, 6)),
			mv(sp(0, 6), sp(1, 7), sp(0, 6)),
			mv(sp(1, 7), sp(0, 8), sp(0, 9)),
			mv(sp(0, 7), sp(2, 9), sp(3, 9)),
			mv(sp(2, 9), sp(2, 7), sp(0, 7)),
			mv(sp(2, 7), sp(1, 7), sp(2, 7)),
			mv(sp(1, 7), sp(2, 8), sp(1, 7)),
			mv(sp(2, 8), sp(2, 9), sp(2, 8)),
			mv(sp(2, 9), sp(1, 9), sp(2, 9)),
			mv(sp(1, 9), sp(1, 8), sp(1, 9)),
		},
	}
	expectedMoves := len(moves[1]) + len(moves[2])

	g := getGame(moves)
	err := g.Play()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if len(g.Events()) != expectedMoves {
		t.Errorf("Expected %d moves, got %d", expectedMoves, len(g.Events()))
	}

	places := g.Places()
	if places[0].Player.ID != g.players[0].ID || places[0].Rank != 1 {
		t.Errorf("Got incorrect places, first player: %+v", places[0])
	}

	if places[1].Player.ID != g.players[1].ID || places[1].Rank != 2 {
		t.Errorf("Got incorrect places, second player: %+v", places[1])
	}
}
