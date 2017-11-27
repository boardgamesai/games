package hearts

import (
	"reflect"
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

func getCards(cards []string) []card.Card {
	hand := []card.Card{}
	for _, c := range cards {
		hand = append(hand, card.FromString(c))
	}
	return hand
}
