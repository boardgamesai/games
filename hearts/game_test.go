package hearts

import (
	"testing"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts/card"
)

func getGame(hands map[int][]string) *Game {
	g := New()
	g.Comms = &CommsMock{
		hands: hands,
	}

	for i := 1; i <= g.NumPlayers(); i++ {
		player := Player{
			Player: game.Player{
				Order: i,
			},
			Runnable: &game.RunnablePlayerMock{},
			Hand:     getHand(hands[i]),
		}
		g.players = append(g.players, &player)
	}

	return g
}

func getHand(cards []string) Hand {
	hand := Hand{}
	for _, c := range cards {
		hand.Add(card.FromString(c))
	}
	return hand
}

func TestPassCards(t *testing.T) {
	tests := []struct {
		hands         map[int][]string
		passDirection PassDirection
		expectedHands map[int][]string
	}{
		{
			map[int][]string{
				1: []string{"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"},
				2: []string{"2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: []string{"2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: []string{"2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
			PassLeft,
			map[int][]string{
				1: []string{"5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC", "2H", "3H", "4H"},
				2: []string{"2C", "3C", "4C", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: []string{"2D", "3D", "4D", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: []string{"5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH", "2S", "3S", "4S"},
			},
		},
		{
			map[int][]string{
				1: []string{"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"},
				2: []string{"2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: []string{"2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: []string{"2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
			PassRight,
			map[int][]string{
				1: []string{"5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC", "2D", "3D", "4D"},
				2: []string{"5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD", "2S", "3S", "4S"},
				3: []string{"2H", "3H", "4H", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: []string{"2C", "3C", "4C", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
		},
		{
			map[int][]string{
				1: []string{"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"},
				2: []string{"2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: []string{"2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: []string{"2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
			PassAcross,
			map[int][]string{
				1: []string{"5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC", "2S", "3S", "4S"},
				2: []string{"5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD", "2H", "3H", "4H"},
				3: []string{"2C", "3C", "4C", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: []string{"2D", "3D", "4D", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
		},
	}

	for _, test := range tests {
		g := getGame(test.hands)

		err := g.passCards(test.passDirection)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		for _, player := range g.players {
			h := getHand(test.expectedHands[player.Order])
			for i := 1; i < len(player.Hand); i++ {
				if player.Hand[i] != h[i] {
					t.Errorf("expected hand: %s got: %s for player %s", h, player.Hand, player)
					break
				}
			}
		}
	}
}

func TestPlayRound(t *testing.T) {
	tests := []struct {
		hands          map[int][]string
		expectedScores []int
	}{
		{
			map[int][]string{
				1: []string{"QC", "QS", "KS", "4D", "2D", "KD", "7H", "6H", "7S", "6S", "5H", "5S", "4S"},
				2: []string{"KC", "4C", "TC", "TD", "7D", "AD", "7C", "9H", "8C", "AC", "4H", "KH", "3C"},
				3: []string{"2C", "5C", "6C", "2S", "JC", "TS", "9C", "QH", "AS", "9S", "TH", "2H", "JH"},
				4: []string{"JS", "8S", "9D", "8D", "6D", "5D", "3D", "8H", "3S", "QD", "3H", "AH", "JD"},
			},
			[]int{0, 0, 22, -6},
		},
		{
			map[int][]string{
				1: []string{"4C", "7C", "TC", "4D", "7D", "TD", "4S", "7S", "TS", "4H", "7H", "TH", "JS"},
				2: []string{"AC", "KC", "QC", "AD", "KD", "QD", "AS", "KS", "QS", "AH", "KH", "QH", "JH"},
				3: []string{"2C", "5C", "8C", "2D", "5D", "8D", "2S", "5S", "8S", "2H", "5H", "8H", "JC"},
				4: []string{"3C", "6C", "9C", "3D", "6D", "9D", "3S", "6S", "9S", "3H", "6H", "9H", "JD"},
			},
			[]int{26, -10, 26, 26},
		},
		{
			map[int][]string{
				1: []string{"4C", "7C", "TC", "4D", "7D", "TD", "4S", "7S", "TS", "4H", "7H", "TH", "JS"},
				2: []string{"AC", "KC", "QC", "JD", "KD", "QD", "AS", "KS", "QS", "AH", "KH", "QH", "JH"},
				3: []string{"2C", "5C", "8C", "AD", "5D", "8D", "2S", "5S", "8S", "2H", "5H", "8H", "JC"},
				4: []string{"3C", "6C", "9C", "3D", "6D", "9D", "3S", "6S", "9S", "3H", "6H", "9H", "2D"},
			},
			[]int{26, 0, 16, 26},
		},
	}

	for _, test := range tests {
		g := getGame(test.hands)

		err := g.playRound()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		for i, score := range test.expectedScores {
			actualScore := g.scores.Totals[g.players[i]]
			if actualScore != score {
				t.Errorf("Got score %d, expected %d", actualScore, score)
			}
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
