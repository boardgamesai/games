package hearts

import (
	"fmt"
	"strings"

	"github.com/boardgamesai/games/game"
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
	ID game.PlayerID
	Hand
}

func (e EventDeal) String() string {
	return fmt.Sprintf("P%d dealt %s", e.ID, e.Hand)
}

type EventPass struct {
	FromID game.PlayerID
	ToID   game.PlayerID
	Cards  []card.Card
}

func (e EventPass) String() string {
	return fmt.Sprintf("P%d passes %s to P%d", e.FromID, e.Cards, e.ToID)
}

type EventPlay struct {
	ID   game.PlayerID
	Card card.Card
}

func (e EventPlay) String() string {
	return fmt.Sprintf("P%d plays %s", e.ID, e.Card)
}

type EventScoreTrick struct {
	ID    game.PlayerID
	Score int
}

func (e EventScoreTrick) String() string {
	return fmt.Sprintf("P%d wins trick, score %d", e.ID, e.Score)
}

type EventScoreRound struct {
	RoundScores map[game.PlayerID]int
	TotalScores map[game.PlayerID]int
}

func (e EventScoreRound) String() string {
	roundVals := []string{}
	totalVals := []string{}

	for playerID, _ := range e.RoundScores {
		score, ok := e.RoundScores[playerID]
		if !ok {
			score = 0
		}

		roundVals = append(roundVals, fmt.Sprintf("P%d:%d", playerID, score))
		totalVals = append(totalVals, fmt.Sprintf("P%d:%d", playerID, e.TotalScores[playerID]))
	}

	return fmt.Sprintf("Round scores: [%s] Total scores: [%s]", strings.Join(roundVals, " "), strings.Join(totalVals, " "))
}
