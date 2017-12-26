package hearts

import (
	"fmt"
	"sort"

	"github.com/boardgamesai/games/game"
)

type Scores struct {
	Rounds []map[*Player]int
	Totals map[*Player]int
}

func NewScores() *Scores {
	scores := Scores{
		Rounds: []map[*Player]int{},
		Totals: map[*Player]int{},
	}
	return &scores
}

func (s *Scores) AddRound(round map[*Player]int) {
	for player, score := range round {
		s.Totals[player] += score
	}
	s.Rounds = append(s.Rounds, round)
}

func (s *Scores) Places() []game.Place {
	type Pair struct {
		score  int
		player *Player
	}
	pairs := []Pair{}

	ties := map[int]int{}
	players := s.Players()

	for _, player := range players {
		score := s.Totals[player]
		pairs = append(pairs, Pair{score: score, player: player})
		ties[score]++
	}

	// We do stable so that game order is maintained among tied players.
	sort.SliceStable(pairs, func(i, j int) bool { return pairs[i].score < pairs[j].score })

	places := []game.Place{}
	rank := 0
	prevScore := -1000000 // Not a plausible score

	for _, pair := range pairs {
		tie := false
		if ties[pair.score] > 1 { // 2+ denotes ties
			tie = true
			if prevScore != pair.score {
				rank++
			}
		} else { // This score is not a tie
			if ties[prevScore] > 1 {
				rank += ties[prevScore] // Skip ahead the number of ties at the previous score
			} else {
				rank++
			}
		}

		place := game.Place{
			Player:   pair.player.Player,
			Rank:     rank,
			Tie:      tie,
			Score:    pair.score,
			HasScore: true,
		}
		places = append(places, place)

		prevScore = pair.score
	}

	return places
}

// Players maintains the list of players in game order. Our map will return them in random order.
func (s *Scores) Players() []*Player {
	players := []*Player{}
	for player, _ := range s.Totals {
		players = append(players, player)
	}
	sort.Slice(players, func(i, j int) bool { return players[i].Order < players[j].Order })
	return players
}

func (s *Scores) String() string {
	players := s.Players()

	str := "    P1  P2  P3  P4\n"
	str += "    --  --  --  --\n"

	for i, round := range s.Rounds {
		str += fmt.Sprintf("%2d", i+1)
		for _, player := range players {
			str += fmt.Sprintf(" %3d", round[player])
		}
		str += "\n"
	}

	str += "    --  --  --  --\n"
	str += "  "
	for _, player := range players {
		str += fmt.Sprintf(" %3d", s.Totals[player])
	}

	return str
}
