package hearts

import (
	"fmt"
	"testing"
)

func getPlayers() []*Player {
	players := make([]*Player, 4)
	for i := 1; i <= 4; i++ {
		player := Player{}
		player.Position = i
		player.Name = fmt.Sprintf("player%d", i)
		players[i-1] = &player
	}

	return players
}

func getRound(scores []int, players []*Player) map[*Player]int {
	round := map[*Player]int{}
	for i := 0; i < 4; i++ {
		round[players[i]] = scores[i]
	}
	return round
}

func TestScores(t *testing.T) {
	players := getPlayers()
	scores := NewScores()

	for player, total := range scores.Totals {
		if total != 0 {
			t.Errorf("On init, didn't get 0 total for player: %v total: %d", player, total)
		}
	}

	tests := []struct {
		round          []int
		expectedTotals []int
	}{
		{[]int{5, 0, -3, 24}, []int{5, 0, -3, 24}},
		{[]int{1, 12, 13, -10}, []int{6, 12, 10, 14}},
		{[]int{-9, 0, 25, 0}, []int{-3, 12, 35, 14}},
	}

	for testIndex, test := range tests {
		round := getRound(test.round, players)
		scores.AddRound(round)

		for i, player := range players {
			if scores.Totals[player] != test.expectedTotals[i] {
				t.Errorf("Wrong total for player: %v expected: %d got: %d test: %d",
					player, test.expectedTotals[i], scores.Totals[player], testIndex)
			}
		}
	}
}

func TestPlaces(t *testing.T) {
	players := getPlayers()
	scores := NewScores()

	tests := []struct {
		round   []int
		players []int
		ranks   []int
		ties    map[int]bool
	}{
		{[]int{4, 8, 7, -3}, []int{4, 1, 3, 2}, []int{1, 2, 3, 4}, map[int]bool{}},
		{[]int{0, 3, 4, 9}, []int{1, 4, 2, 3}, []int{1, 2, 3, 3}, map[int]bool{3: true, 4: true}},
		{[]int{7, 0, 3, 6}, []int{1, 2, 4, 3}, []int{1, 1, 3, 4}, map[int]bool{1: true, 2: true}},
		{[]int{5, 5, 2, 4}, []int{1, 2, 3, 4}, []int{1, 1, 1, 1}, map[int]bool{1: true, 2: true, 3: true, 4: true}},
		{[]int{5, 5, 1, 5}, []int{3, 1, 2, 4}, []int{1, 2, 2, 2}, map[int]bool{2: true, 3: true, 4: true}},
		{[]int{1, 1, 13, 1}, []int{1, 2, 4, 3}, []int{1, 1, 1, 4}, map[int]bool{1: true, 2: true, 3: true}},
		{[]int{4, 1, 7, 4}, []int{2, 1, 4, 3}, []int{1, 2, 2, 4}, map[int]bool{2: true, 3: true}},
		{[]int{5, -10, 11, 0}, []int{2, 4, 1, 3}, []int{1, 2, 3, 4}, map[int]bool{}},
	}

	for _, test := range tests {
		round := getRound(test.round, players)
		scores.AddRound(round)
		places := scores.Places()

		for i, place := range places {
			if place.Player.Name != fmt.Sprintf("player%d", test.players[i]) {
				t.Errorf("Wrong player: %s expected: player%d place: %+v", place.Player.Name, test.players[i], place)
			}

			if place.Rank != test.ranks[i] {
				t.Errorf("Wrong rank: %d expected: %d place: %+v", place.Rank, test.ranks[i], place)
			}

			if place.Tie && !test.ties[place.Rank] {
				t.Errorf("Got unexpected tie: %+v", place)
			} else if !place.Tie && test.ties[place.Rank] {
				t.Errorf("Expected tie but didn't get one, place: %+v", place)
			}
		}
	}
}
