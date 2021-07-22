package hearts

import (
	"encoding/json"
	"fmt"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/hearts/card"
	"github.com/boardgamesai/games/util"
)

type Game struct {
	game.Game
	Comms   AIComms
	players []*Player
	deck    *Deck
	scores  *Scores
}

func New() *Game {
	g := Game{
		players: []*Player{},
	}
	g.Name = game.Hearts

	for i := 0; i < g.MetaData().NumPlayers; i++ {
		p := Player{
			Player: game.Player{},
		}
		g.players = append(g.players, &p)
	}

	g.reset()
	return &g
}

func (g *Game) Play() error {
	// Wipe out any previous state
	g.reset()
	g.shufflePlayers()

	// We need to write down our setup
	setupEvent := EventSetup{
		Players: []EventSetupPlayer{},
	}

	// Launch the player processes
	for _, player := range g.players {
		defer player.CleanUp()
		defer g.SetOutput(player.ID, player)

		err := player.Run()
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}

		err = g.Comms.Setup(player, g.players)
		if err != nil {
			return fmt.Errorf("player %s failed to setup, err: %s", player, err)
		}

		// Keep track of our setup
		esp := EventSetupPlayer{
			ID:       player.ID,
			Position: player.Position,
		}
		setupEvent.Players = append(setupEvent.Players, esp)
	}

	g.EventLog.AddNone(setupEvent)

	passDirection := PassLeft
	for !g.gameOver() {
		err := g.dealCards()
		if err != nil {
			// TODO how to correctly handle when we bomb out here? Who wins?
			return err
		}

		if passDirection != PassNone {
			err = g.passCards(passDirection)
			if err != nil {
				// TODO how to correctly handle when we bomb out here? Who wins?
				return err
			}
		}

		err = g.playRound()
		if err != nil {
			// TODO how to correctly handle when we bomb out here? Who wins?
			return err
		}

		passDirection = passDirection.Next()
	}

	// We're done, set the places for each player.
	g.SetPlaces(g.scores.Places())

	return nil
}

func (g *Game) Players() []*game.Player {
	players := []*game.Player{}
	for _, p := range g.players {
		players = append(players, &(p.Player))
	}
	return players
}

func (g *Game) Events() []fmt.Stringer {
	events := make([]fmt.Stringer, len(g.EventLog))

	for i, event := range g.EventLog {
		var eStr fmt.Stringer

		switch event.Type {
		case EventTypeSetup:
			e := EventSetup{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		case EventTypeDeal:
			e := EventDeal{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		case EventTypePass:
			e := EventPass{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		case EventTypePlay:
			e := EventPlay{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		case EventTypeScoreTrick:
			e := EventScoreTrick{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		case EventTypeScoreRound:
			e := EventScoreRound{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		}

		events[i] = eStr
	}

	return events
}

func (g *Game) reset() {
	g.Game.Reset()
	g.deck = NewDeck()
	g.scores = NewScores()
	if g.Comms == nil {
		g.Comms = NewComms(g)
	}
}

func (g *Game) dealCards() error {
	g.deck.Shuffle()
	for _, player := range g.players {
		player.Hand = Hand{}
		for i := 0; i < 13; i++ {
			player.Hand.Add(g.deck.DealCard())
		}
		player.Hand.Sort()

		e := EventDeal{
			ID:   player.ID,
			Hand: player.Hand,
		}
		g.EventLog.Add(e, []game.PlayerID{player.ID})
	}

	return nil
}

func (g *Game) passCards(passDirection PassDirection) error {
	// First collect all the passes...
	passes := map[*Player]PassMove{}
	for _, player := range g.players {
		passMove, err := g.Comms.GetPassMove(player, passDirection)
		if err != nil {
			// TODO: who wins if this happens?
			return fmt.Errorf("player %s failed to pass cards, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.isValidPass(player, passMove)
		if err != nil {
			// TODO: who wins if this happens?
			return fmt.Errorf("player %s made invalid pass: %+v err: %s", player, passMove, err)
		}

		passes[player] = passMove
	}

	// ... and now distribute the passes. We do this so that no player gets their passed cards before
	// they choose which to pass.
	for passer, passMove := range passes {
		recipient := g.getPassRecipient(passer, passDirection)
		for _, card := range passMove.Cards {
			passer.Hand.Remove(card) // Remove the card from the passer's hand,
			recipient.Hand.Add(card) // and add it to the recipient's hand.
		}
		recipient.Hand.Sort()

		e := EventPass{
			FromID: passer.ID,
			ToID:   recipient.ID,
			Cards:  passMove.Cards,
		}
		g.EventLog.Add(e, []game.PlayerID{passer.ID, recipient.ID})
	}

	return nil
}

func (g *Game) isValidPass(p *Player, m PassMove) error {
	if len(m.Cards) != 3 {
		return fmt.Errorf("player %s passed %d cards, expected 3", p, len(m.Cards))
	}

	// Make sure each card is actually in their hand
	for _, passCard := range m.Cards {
		if !p.Hand.Contains(passCard) {
			return fmt.Errorf("player %s passed card not in their hand: %s", p, passCard)
		}
	}

	return nil
}

func (g *Game) getPassRecipient(p *Player, passDirection PassDirection) *Player {
	playerIndex := -1
	for i, player := range g.players {
		if player == p {
			playerIndex = i
			break
		}
	}

	addon := 0
	switch passDirection {
	case PassLeft:
		addon = 1
	case PassAcross:
		addon = 2
	case PassRight:
		addon = 3
	}

	return g.players[(playerIndex+addon)%4]
}

func (g *Game) playRound() error {
	// To kick off the round, we need to know who has the two of clubs.
	turn := -1

	for i, player := range g.players {
		for _, c := range player.Hand {
			if c.Suit == card.Clubs && c.Rank == card.Two {
				turn = i
				break
			}
		}
	}

	scores := map[*Player]int{}
	score := 0
	heartsBroken := false
	tookPoints := map[*Player]bool{}
	var err error

	for i := 0; i < 13; i++ {
		turn, score, err = g.playTrick(turn, i, heartsBroken)
		if err != nil {
			return err
		}

		// We deduce whether a heart got played or not based on the score.
		// The only scores where no hearts were played are 0, 13 (QS only), -10 (JD only).
		// Any other score means a heart was in the mix.
		if !heartsBroken && score != 0 && score != 13 && score != -10 {
			heartsBroken = true
		}

		scores[g.players[turn]] += score

		if score != 0 && score != -10 {
			tookPoints[g.players[turn]] = true
		}
	}

	// Before this round's in the books, check for a moonshot, which would change everything.
	if len(tookPoints) == 1 {
		// Only one player taking points: that's a moonshot
		moonshotter := &Player{}
		for player := range tookPoints {
			moonshotter = player
		}

		for _, player := range g.players {
			if player == moonshotter {
				scores[player] -= 26
			} else {
				scores[player] += 26
			}
		}
	}

	g.scores.AddRound(scores)

	eventScores := map[game.PlayerID]int{}
	for player, score := range scores {
		eventScores[player.ID] = score
	}

	totalScores := map[game.PlayerID]int{}
	for player, score := range g.scores.Totals {
		totalScores[player.ID] = score
	}

	e := EventScoreRound{
		RoundScores: eventScores,
		TotalScores: totalScores,
	}
	g.EventLog.AddAll(e)

	return nil
}

func (g *Game) playTrick(turn int, trickCount int, heartsBroken bool) (int, int, error) {
	trick := []card.Card{}
	plays := map[card.Card]game.PlayerID{}
	turns := map[card.Card]int{}

	// Collect a play from each player
	for i := 0; i < 4; i++ {
		player := g.players[turn]
		move, err := g.Comms.GetPlayMove(player, trick)
		if err != nil {
			return -1, -1, err
		}

		err = g.isValidPlay(player, move, trick, trickCount, heartsBroken)
		if err != nil {
			return -1, -1, err
		}

		trick = append(trick, move.Card)
		player.Hand.Remove(move.Card)
		plays[move.Card] = player.ID
		turns[move.Card] = turn
		turn = util.Increment(turn, 0, 3)

		e := EventPlay{
			ID:   player.ID,
			Card: move.Card,
		}
		g.EventLog.AddAll(e)
	}

	// Now see what the trick is worth and who gets it.
	topCard, score := g.evaluateTrick(trick)

	e := EventScoreTrick{
		ID:    plays[topCard],
		Score: score,
	}
	g.EventLog.AddAll(e)

	return turns[topCard], score, nil
}

func (g *Game) isValidPlay(p *Player, m PlayMove, trick []card.Card, trickCount int, heartsBroken bool) error {
	// First make sure the card is actually in their hand.
	if !p.Hand.Contains(m.Card) {
		return fmt.Errorf("player %s played card not in their hand: %s", p, m.Card)
	}

	// Now see if the card can actually be played.
	valid := false
	for _, card := range p.Hand.PossiblePlays(trick, trickCount, heartsBroken) {
		if m.Card == card {
			valid = true
		}
	}
	if !valid {
		return fmt.Errorf("player %s played invalid card: %s", p, m.Card)
	}

	return nil
}

// evaluateTrick returns the winning card and the score of the trick
func (g *Game) evaluateTrick(trick []card.Card) (card.Card, int) {
	winner := 0
	score := 0

	for i, c := range trick {
		if c.Rank == card.Queen && c.Suit == card.Spades {
			score += 13
		} else if c.Rank == card.Jack && c.Suit == card.Diamonds {
			score -= 10
		} else if c.Suit == card.Hearts {
			score++
		}

		if i > 0 {
			if c.Suit == trick[winner].Suit && c.Index() > trick[winner].Index() {
				winner = i
			}
		}
	}

	return trick[winner], score
}

func (g *Game) shufflePlayers() {
	util.Shuffle(g.players)

	for i := 1; i <= 4; i++ {
		g.players[i-1].Position = i
	}
}

func (g *Game) gameOver() bool {
	for _, total := range g.scores.Totals {
		if total >= 100 {
			return true
		}
	}

	return false
}
