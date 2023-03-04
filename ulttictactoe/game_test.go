package ulttictactoe

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

func mv(col, row, subCol, subRow int) Move {
	return Move{
		Col:    col,
		Row:    row,
		SubCol: subCol,
		SubRow: subRow,
	}
}

func TestPlayGame(t *testing.T) {
	games := []map[int][]Move{
		{
			// Shortest possible game
			1: {
				mv(1, 2, 2, 2),
				mv(1, 2, 2, 1),
				mv(1, 2, 2, 0),
				mv(1, 1, 2, 2),
				mv(1, 1, 2, 1),
				mv(1, 1, 2, 0),
				mv(1, 0, 2, 2),
				mv(1, 0, 2, 1),
				mv(1, 0, 2, 0),
			},
			2: {
				mv(2, 2, 1, 2),
				mv(2, 1, 1, 2),
				mv(2, 0, 1, 2),
				mv(2, 2, 1, 1),
				mv(2, 1, 1, 1),
				mv(2, 0, 1, 1),
				mv(2, 2, 1, 0),
				mv(2, 1, 1, 0),
			},
		},
		{
			// Normal game
			1: {
				mv(1, 1, 1, 1),
				mv(2, 0, 1, 2),
				mv(1, 2, 1, 0),
				mv(1, 0, 1, 2),
				mv(2, 2, 2, 2),
				mv(2, 1, 2, 1),
				mv(0, 2, 0, 2),
				mv(1, 0, 2, 0),
				mv(0, 1, 0, 1),
				mv(1, 1, 2, 1),
				mv(2, 2, 1, 1),
				mv(1, 0, 2, 2),
				mv(0, 0, 0, 0),
				mv(2, 1, 1, 2),
				mv(0, 2, 2, 2),
				mv(0, 2, 0, 0),
				mv(0, 2, 0, 1), // Wins 0, 2
				mv(2, 2, 1, 0),
				mv(0, 1, 0, 0),
				mv(0, 1, 2, 0),
				mv(0, 0, 2, 0),
				mv(2, 1, 2, 0),
				mv(1, 0, 2, 1), // Wins 1, 0
				mv(0, 0, 1, 0), // Wins 0, 0
				mv(0, 1, 1, 0), // Wins 0, 1
			},
			2: {
				mv(1, 1, 2, 0),
				mv(1, 2, 1, 2),
				mv(1, 0, 1, 0),
				mv(1, 2, 2, 2),
				mv(2, 2, 2, 1),
				mv(2, 1, 0, 2),
				mv(0, 2, 1, 0),
				mv(2, 0, 0, 1),
				mv(0, 1, 1, 1),
				mv(2, 1, 2, 2),
				mv(1, 1, 1, 0),
				mv(2, 2, 0, 0),
				mv(0, 0, 2, 1),
				mv(1, 2, 0, 2), // Wins 1,2
				mv(2, 2, 0, 2),
				mv(0, 0, 0, 2),
				mv(0, 1, 2, 2),
				mv(1, 0, 0, 1),
				mv(0, 0, 0, 1),
				mv(2, 0, 0, 0),
				mv(2, 0, 2, 1),
				mv(2, 0, 1, 0),
				mv(2, 1, 0, 0),
				mv(0, 1, 0, 2),
			},
		},
	}
	for _, moves := range games {
		expectedMoves := len(moves[1]) + len(moves[2]) + 1 // +1 for the setup event

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
}
