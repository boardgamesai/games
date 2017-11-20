package hearts

import (
	"fmt"
	"strings"

	"github.com/boardgamesai/games/hearts/card"
)

const (
	EventTypeDeal       = "deal"
	EventTypePass       = "pass"
	EventTypePlay       = "play"
	EventTypeScoreTrick = "scoretrick"
	EventTypeScoreRound = "scoreround"
)

type EventDeal struct {
	Order int
	Hand
}

func (e EventDeal) String() string {
	return fmt.Sprintf("P%d dealt %s", e.Order, e.Hand)
}

type EventPass struct {
	FromOrder int
	ToOrder   int
	Cards     []card.Card
}

func (e EventPass) String() string {
	return fmt.Sprintf("P%d passes %s to P%d", e.FromOrder, e.Cards, e.ToOrder)
}

type EventPlay struct {
	Order int
	Card  card.Card
}

func (e EventPlay) String() string {
	return fmt.Sprintf("P%d plays %s", e.Order, e.Card)
}

type EventScoreTrick struct {
	Order int
	Score int
}

func (e EventScoreTrick) String() string {
	return fmt.Sprintf("P%d wins trick, score %d", e.Order, e.Score)
}

type EventScoreRound struct {
	RoundScores map[int]int
	TotalScores map[int]int
}

func (e EventScoreRound) String() string {
	roundVals := make([]string, 4)
	totalVals := make([]string, 4)

	for i := 1; i <= 4; i++ {
		score, ok := e.RoundScores[i]
		if !ok {
			score = 0
		}
		roundVals[i-1] = fmt.Sprintf("P%d:%d", i, score)
		totalVals[i-1] = fmt.Sprintf("P%d:%d", i, e.TotalScores[i])
	}

	return fmt.Sprintf("Round scores: [%s] Total scores: [%s]", strings.Join(roundVals, " "), strings.Join(totalVals, " "))
}
