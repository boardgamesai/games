package hearts

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts/card"
)

func getGame(hands map[int][]string) *Game {
	g := New()

	for i := 0; i < g.MetaData().NumPlayers; i++ {
		g.players[i].Player.ID = game.PlayerID(i + 1)
		g.players[i].Position = i + 1
		g.players[i].Player.Name = fmt.Sprintf("player%d", i+1)
		g.players[i].Player.Runnable = &game.RunnablePlayerMock{}
		g.players[i].Hand = getHand(hands[i+1])
	}

	g.Comms = &CommsMock{
		hands: hands,
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
				1: {"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"},
				2: {"2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: {"2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: {"2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
			PassLeft,
			map[int][]string{
				1: {"5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC", "2H", "3H", "4H"},
				2: {"2C", "3C", "4C", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: {"2D", "3D", "4D", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: {"2S", "3S", "4S", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
		},
		{
			map[int][]string{
				1: {"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"},
				2: {"2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: {"2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: {"2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
			PassRight,
			map[int][]string{
				1: {"5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC", "2D", "3D", "4D"},
				2: {"5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD", "2S", "3S", "4S"},
				3: {"5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS", "2H", "3H", "4H"},
				4: {"2C", "3C", "4C", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
		},
		{
			map[int][]string{
				1: {"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"},
				2: {"2D", "3D", "4D", "5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD"},
				3: {"2S", "3S", "4S", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: {"2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
			},
			PassAcross,
			map[int][]string{
				1: {"5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC", "2S", "3S", "4S"},
				2: {"5D", "6D", "7D", "8D", "9D", "TD", "JD", "QD", "KD", "AD", "2H", "3H", "4H"},
				3: {"2C", "3C", "4C", "5S", "6S", "7S", "8S", "9S", "TS", "JS", "QS", "KS", "AS"},
				4: {"2D", "3D", "4D", "5H", "6H", "7H", "8H", "9H", "TH", "JH", "QH", "KH", "AH"},
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
			h := getHand(test.expectedHands[player.Position])
			for i := 1; i < len(player.Hand); i++ {
				if player.Hand[i] != h[i] {
					t.Errorf("expected hand: %s got: %s for player %s", h, player.Hand, player)
					break
				}
			}
		}
	}
}

func TestIsValidPass(t *testing.T) {
	tests := []struct {
		pass     []string
		expected error
	}{
		{
			[]string{"2C", "3C", "4C"},
			nil,
		},
		{
			[]string{"2C", "3C"},
			InvalidPassError{},
		},
		{
			[]string{"2C", "3C", "4C", "5C"},
			InvalidPassError{},
		},
		{
			[]string{"2C", "3C", "4H"},
			InvalidPassError{},
		},
		{
			[]string{"2C", "3C", "3C"},
			InvalidPassError{},
		},
	}

	g := New()
	for _, test := range tests {
		move := PassMove{
			Cards: getCards(test.pass),
		}
		hand := getHand([]string{"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"})
		err := g.isValidPass(hand, move)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expected) {
			t.Errorf("expected err: %s got: %s", test.expected, err)
		}
	}
}

func TestIsValidPlay(t *testing.T) {
	tests := []struct {
		trick    []string
		play     string
		expected error
	}{
		{
			[]string{"5D"},
			"4D",
			nil,
		},
		{
			[]string{"5D"},
			"3S", // Not in hand
			InvalidPlayError{},
		},
		{
			[]string{"5D"},
			"KS", // Can't be played
			IllegalPlayError{},
		},
	}

	g := New()
	for _, test := range tests {
		move := PlayMove{
			Card: card.FromString(test.play),
		}
		hand := getHand([]string{"KS", "4D", "7C", "TH", "JH", "QH", "KH"})
		err := g.isValidPlay(hand, move, getCards(test.trick), 6, false)
		if reflect.TypeOf(err) != reflect.TypeOf(test.expected) {
			t.Errorf("expected err: %s got: %s", test.expected, err)
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
				1: {"QC", "QS", "KS", "4D", "2D", "KD", "7H", "6H", "7S", "6S", "5H", "5S", "4S"},
				2: {"KC", "4C", "TC", "TD", "7D", "AD", "7C", "9H", "8C", "AC", "4H", "KH", "3C"},
				3: {"2C", "5C", "6C", "2S", "JC", "TS", "9C", "QH", "AS", "9S", "TH", "2H", "JH"},
				4: {"JS", "8S", "9D", "8D", "6D", "5D", "3D", "8H", "3S", "QD", "3H", "AH", "JD"},
			},
			[]int{0, 0, 22, -6},
		},
		{
			map[int][]string{
				1: {"4C", "7C", "TC", "4D", "7D", "TD", "4S", "7S", "TS", "4H", "7H", "TH", "JS"},
				2: {"AC", "KC", "QC", "AD", "KD", "QD", "AS", "KS", "QS", "AH", "KH", "QH", "JH"},
				3: {"2C", "5C", "8C", "2D", "5D", "8D", "2S", "5S", "8S", "2H", "5H", "8H", "JC"},
				4: {"3C", "6C", "9C", "3D", "6D", "9D", "3S", "6S", "9S", "3H", "6H", "9H", "JD"},
			},
			[]int{26, -10, 26, 26},
		},
		{
			map[int][]string{
				1: {"4C", "7C", "TC", "4D", "7D", "TD", "4S", "7S", "TS", "4H", "7H", "TH", "JS"},
				2: {"AC", "KC", "QC", "JD", "KD", "QD", "AS", "KS", "QS", "AH", "KH", "QH", "JH"},
				3: {"2C", "5C", "8C", "AD", "5D", "8D", "2S", "5S", "8S", "2H", "5H", "8H", "JC"},
				4: {"3C", "6C", "9C", "3D", "6D", "9D", "3S", "6S", "9S", "3H", "6H", "9H", "2D"},
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

func TestEvaluateTrick(t *testing.T) {
	tests := []struct {
		trick     []string
		expWinner string
		expScore  int
	}{
		{[]string{"2C", "3C", "4C", "5C"}, "5C", 0},
		{[]string{"2C", "3C", "5C", "4C"}, "5C", 0},
		{[]string{"2C", "5C", "3C", "4C"}, "5C", 0},
		{[]string{"5C", "2C", "3C", "4C"}, "5C", 0},
		{[]string{"3C", "9D", "TH", "JS"}, "3C", 1},
		{[]string{"7H", "4S", "8H", "AC"}, "8H", 2},
		{[]string{"7H", "AH", "8H", "AC"}, "AH", 3},
		{[]string{"7H", "KH", "8H", "AH"}, "AH", 4},
		{[]string{"8D", "2D", "3D", "QS"}, "8D", 13},
		{[]string{"8D", "2H", "3D", "QS"}, "8D", 14},
		{[]string{"8D", "2H", "3H", "QS"}, "8D", 15},
		{[]string{"8H", "2H", "3H", "QS"}, "8H", 16},
		{[]string{"8D", "2D", "JD", "QS"}, "JD", 3},
		{[]string{"8D", "2H", "JD", "QS"}, "JD", 4},
		{[]string{"8H", "2H", "JD", "QS"}, "8H", 5},
		{[]string{"8D", "QD", "JD", "7C"}, "QD", -10},
		{[]string{"8D", "QH", "JD", "7C"}, "JD", -9},
		{[]string{"8D", "QH", "JD", "7H"}, "JD", -8},
		{[]string{"8H", "QH", "JD", "7H"}, "QH", -7},
		{[]string{"5C", "2D", "KS", "8D"}, "5C", 0},
	}

	g := New()

	for _, test := range tests {
		topCard, score := g.evaluateTrick(getCards(test.trick))
		if topCard != card.FromString(test.expWinner) {
			t.Errorf("Got topCard %s for trick %+v, expected %s", topCard, test.trick, test.expWinner)
		}
		if score != test.expScore {
			t.Errorf("Got score %d for trick %+v, expected %d", score, test.trick, test.expScore)
		}
	}
}
