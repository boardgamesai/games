package hearts

import (
	"testing"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts/card"
)

func getGame(hands map[int][]string) *Game {
	g := New()

	players := []*Player{}
	for i := 1; i <= g.NumPlayers(); i++ {
		hand := Hand{}
		for _, c := range hands[i] {
			hand.Add(card.FromString(c))
		}

		player := Player{
			Player: game.Player{
				Order: i,
			},
			Runnable: &game.RunnablePlayerMock{},
			Hand:     hand,
		}
		players = append(players, &player)
	}
	g.players = players
	g.Comms = NewCommsMock(hands)

	return g
}

func TestPlayRound(t *testing.T) {
	hands := map[int][]string{
		1: []string{"QC", "QS", "KS", "4D", "2D", "KD", "7H", "6H", "7S", "6S", "5H", "5S", "4S"},
		2: []string{"KC", "4C", "TC", "TD", "7D", "AD", "7C", "9H", "8C", "AC", "4H", "KH", "3C"},
		3: []string{"2C", "5C", "6C", "2S", "JC", "TS", "9C", "QH", "AS", "9S", "TH", "2H", "JH"},
		4: []string{"JS", "8S", "9D", "8D", "6D", "5D", "3D", "8H", "3S", "QD", "3H", "AH", "JD"},
	}
	g := getGame(hands)

	err := g.playRound()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	expectedScores := []int{0, 0, 22, -6}
	for i, score := range expectedScores {
		actualScore := g.scores.Totals[g.players[i]]
		if actualScore != score {
			t.Errorf("Got score %d, expected %d", actualScore, score)
		}
	}
}

func getTrick(cards []string) []card.Card {
	trick := []card.Card{}
	for _, c := range cards {
		trick = append(trick, card.FromString(c))
	}
	return trick
}

func TestEvaluateTrick(t *testing.T) {
	tests := []struct {
		trick     []card.Card
		expWinner int
		expScore  int
	}{
		{getTrick([]string{"2C", "3C", "4C", "5C"}), 3, 0},
		{getTrick([]string{"2C", "3C", "5C", "4C"}), 2, 0},
		{getTrick([]string{"2C", "5C", "3C", "4C"}), 1, 0},
		{getTrick([]string{"5C", "2C", "3C", "4C"}), 0, 0},
		{getTrick([]string{"3C", "9D", "TH", "JS"}), 0, 1},
		{getTrick([]string{"7H", "4S", "8H", "AC"}), 2, 2},
		{getTrick([]string{"7H", "AH", "8H", "AC"}), 1, 3},
		{getTrick([]string{"7H", "KH", "8H", "AH"}), 3, 4},
		{getTrick([]string{"8D", "2D", "3D", "QS"}), 0, 13},
		{getTrick([]string{"8D", "2H", "3D", "QS"}), 0, 14},
		{getTrick([]string{"8D", "2H", "3H", "QS"}), 0, 15},
		{getTrick([]string{"8H", "2H", "3H", "QS"}), 0, 16},
		{getTrick([]string{"8D", "2D", "JD", "QS"}), 2, 3},
		{getTrick([]string{"8D", "2H", "JD", "QS"}), 2, 4},
		{getTrick([]string{"8H", "2H", "JD", "QS"}), 0, 5},
		{getTrick([]string{"8D", "QD", "JD", "7C"}), 1, -10},
		{getTrick([]string{"8D", "QH", "JD", "7C"}), 2, -9},
		{getTrick([]string{"8D", "QH", "JD", "7H"}), 2, -8},
		{getTrick([]string{"8H", "QH", "JD", "7H"}), 1, -7},
		{getTrick([]string{"5C", "2D", "KS", "8D"}), 0, 0},
	}

	g := New()

	for _, test := range tests {
		winner, score := g.evaluateTrick(test.trick)
		if winner != test.expWinner {
			t.Errorf("Got winner %d for trick %+v, expected %d", winner, test.trick, test.expWinner)
		}
		if score != test.expScore {
			t.Errorf("Got score %d for trick %+v, expected %d", score, test.trick, test.expScore)
		}
	}
}
