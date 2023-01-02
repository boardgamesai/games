package hearts

import (
	"reflect"
	"strings"
	"testing"

	"github.com/boardgamesai/games/hearts/card"
)

func TestAddRemoveContains(t *testing.T) {
	h := Hand{}
	h.Add(card.FromString("5H"))
	h.Add(card.FromString("JC"))

	if len(h) != 2 {
		t.Errorf("Hand is not length 2")
	}

	ok := h.Contains(card.FromString("5H"))
	if !ok {
		t.Errorf("Hand did not contain 5H")
	}

	ok = h.Contains(card.FromString("AD"))
	if ok {
		t.Errorf("Hand contained AD but should not have")
	}

	ok = h.Remove(card.FromString("5H"))
	if !ok {
		t.Errorf("Couldn't remove 5H")
	}

	if len(h) != 1 {
		t.Errorf("Hand is not length 1")
	}

	ok = h.Contains(card.FromString("JC"))
	if !ok {
		t.Errorf("Hand did not contain JC")
	}

	ok = h.Remove(card.FromString("JC"))
	if !ok {
		t.Errorf("Couldn't remove JC")
	}

	if len(h) != 0 {
		t.Errorf("Hand is not length 0")
	}
}

func TestPossiblePlays(t *testing.T) {
	tests := []struct {
		hand          []string
		trick         []string
		trickCount    int
		heartsBroken  bool
		expectedPlays []string
	}{
		{
			[]string{"2C", "5C", "6C", "JC", "4D", "9D", "TD", "JD", "KD", "3S", "AS", "6H", "QH"},
			[]string{},
			0,
			false,
			[]string{"2C"},
		},
		{
			[]string{"3C", "5C", "6C", "JC", "4D", "9D", "TD", "JD", "KD", "3S", "AS", "6H", "QH"},
			[]string{"2C"},
			0,
			false,
			[]string{"3C", "5C", "6C", "JC"},
		},
		{
			[]string{"3C", "5C", "6C", "4D", "9D", "TD", "JD", "KD", "3S", "AS", "6H", "QH"},
			[]string{"8D", "3D", "2D"},
			1,
			false,
			[]string{"4D", "9D", "TD", "JD", "KD"},
		},
		{
			[]string{"3C", "5C", "6C", "4D", "9D", "TD", "JD", "KD", "3S", "AS", "6H", "QH"},
			[]string{"JS", "AC"},
			1,
			false,
			[]string{"3S", "AS"},
		},
		{
			[]string{"3C", "5C", "6C", "4D", "9D", "JD", "KD", "3S", "AS", "6H", "QH"},
			[]string{"2H", "AC"},
			2,
			true,
			[]string{"6H", "QH"},
		},
		{
			[]string{"3C", "5C", "6C", "4D", "9D", "JD", "KD", "6H", "QH"},
			[]string{"9S", "8S"},
			4,
			false,
			[]string{"3C", "5C", "6C", "4D", "9D", "JD", "KD", "6H", "QH"},
		},
		{
			[]string{"4D", "9D", "TD", "JD", "KD", "3S", "4S", "5S", "6S", "JS", "6H", "QH", "KH"},
			[]string{"2C"},
			0,
			false,
			[]string{"4D", "9D", "TD", "JD", "KD", "3S", "4S", "5S", "6S", "JS"},
		},
		{
			[]string{"4D", "9D", "TD", "JD", "KD", "3S", "4S", "5S", "QS", "KS", "6H", "QH", "KH"},
			[]string{"2C"},
			0,
			false,
			[]string{"4D", "9D", "TD", "JD", "KD", "3S", "4S", "5S", "KS"},
		},
		{
			[]string{"5C", "6C", "4D", "9D", "TD", "JD", "KD", "3S", "AS", "6H", "QH"},
			[]string{},
			2,
			false,
			[]string{"5C", "6C", "4D", "9D", "TD", "JD", "KD", "3S", "AS"},
		},
		{
			[]string{"5C", "6C", "4D", "9D", "TD", "JD", "KD", "3S", "AS", "6H", "QH"},
			[]string{},
			2,
			true,
			[]string{"5C", "6C", "4D", "9D", "TD", "JD", "KD", "3S", "AS", "6H", "QH"},
		},
		{
			[]string{"QS", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "QH", "KH", "AH"},
			[]string{"2C"},
			0,
			false,
			[]string{"QS", "2H", "3H", "4H", "5H", "6H", "7H", "8H", "9H", "TH", "QH", "KH", "AH"},
		},
	}

	for _, test := range tests {
		h := Hand{}
		for _, c := range getCards(test.hand) {
			h.Add(c)
		}

		plays := h.PossiblePlays(getCards(test.trick), test.trickCount, test.heartsBroken)
		if !reflect.DeepEqual(plays, getCards(test.expectedPlays)) {
			t.Errorf("Got plays %s, expected %s for hand %s, trick %s, trickCount %d, heartsBroken %t",
				plays, test.expectedPlays, test.hand, test.trick, test.trickCount, test.heartsBroken)
		}
	}
}

func TestIsValidPass(t *testing.T) {
	tests := []struct {
		pass        []string
		expectedErr string
	}{
		{
			[]string{"2C", "3C", "4C"},
			"",
		},
		{
			[]string{"2C", "3C"},
			"didn't get",
		},
		{
			[]string{"2C", "3C", "4C", "5C"},
			"didn't get",
		},
		{
			[]string{"2C", "3C", "4H"},
			"not in hand",
		},
		{
			[]string{"2C", "3C", "3C"},
			"duplicated card",
		},
	}

	for _, test := range tests {
		hand := getHand([]string{"2C", "3C", "4C", "5C", "6C", "7C", "8C", "9C", "TC", "JC", "QC", "KC", "AC"})
		err := hand.IsValidPass(getCards(test.pass))
		if err != nil {
			if test.expectedErr == "" {
				t.Errorf("expected no error, got %s, pass: %+v", err, test.pass)
			} else if !strings.Contains(err.Error(), test.expectedErr) {
				t.Errorf("expected error %s, got %s, pass: %+v", test.expectedErr, err, test.pass)
			}
		} else if test.expectedErr != "" {
			t.Errorf("expected error %s, got none, pass: %+v", test.expectedErr, test.pass)
		}
	}
}

func getCards(cards []string) []card.Card {
	hand := []card.Card{}
	for _, c := range cards {
		hand = append(hand, card.FromString(c))
	}
	return hand
}
