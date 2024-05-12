package liarsdice

import (
	"encoding/json"
	"fmt"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/util"
)

type Game struct {
	game.Game[*Player, *Board, AIComms]
}

func New() *Game {
	g := Game{}
	g.Name = game.LiarsDice
	g.InitPlayers(NewPlayer)
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
	for _, player := range g.Players {
		defer player.CleanUp()
		defer g.SetOutput(player.ID, player)

		err := player.Run()
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}

		err = g.Comms.Setup(player, g.Players)
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
	g.sendRollEvents()

	playerTurn := 0
	for !g.gameOver() {
		player := g.Players[playerTurn]
		move, err := g.Comms.GetMove(player)
		if err != nil {
			g.setLoser(player)
			switch e := err.(type) {
			// If this is a DQError, we need to augment it with the player ID,
			// which we may not know about where the error occurred
			case game.DQError:
				return g.AddDQErrorID(&e, player.ID)
			case *game.DQError:
				return g.AddDQErrorID(e, player.ID)
			}
			return err
		}

		err = g.Board.ApplyMove(move, player)
		if err != nil {
			g.setLoser(player)
			return game.DQError{
				ID:   player.ID,
				Type: game.DQTypeInvalidMove,
				Msg:  err.Error(),
			}
		}

		if !move.Challenge {
			// New bid is simple, just log it
			e := EventMove{
				ID:   player.ID,
				Move: move,
			}
			g.EventLog.AddAll(e)

			// Did dice get shown here? If so we need to send an event about the remaining dice re-roll
			// (the board already did the re-roll)
			if len(move.ShowDice) > 0 {
				e := EventRoll{
					ID:   player.ID,
					Dice: g.Board.DiceHidden[player].Values,
				}
				g.EventLog.Add(e, []game.PlayerID{player.ID})
			}
		} else {
			// Challenge will result in dice changes / potential eliminations
			var eliminated *Player

			changes := map[game.PlayerID]int{}
			for p, change := range g.Board.Outcome.DiceChanges {
				changes[p.ID] = change

				// Did a player just get eliminated?
				if change < 0 && g.Board.DiceHidden[p].Count() == 0 {
					eliminated = p
				}
			}
			e := EventChallenge{
				ID:             player.ID,
				Bid:            g.Board.Outcome.Bid,
				ActualQuantity: g.Board.Outcome.ActualQuantity,
				DiceChange:     changes,
			}
			if eliminated != nil {
				e.Eliminated = eliminated.ID
			}
			g.EventLog.AddAll(e)

			if eliminated != nil {
				place := game.Place{
					Player: eliminated.Player,
					Rank:   g.MetaData().NumPlayers - len(g.Places()),
				}
				g.AddPlace(place)
			}

			if !g.gameOver() {
				g.sendRollEvents()
			}
		}

		playerTurn = util.Increment(playerTurn, 0, g.MetaData().NumPlayers-1)
		for g.Board.DiceHidden[g.Players[playerTurn]].Count() == 0 {
			// Skip over eliminated players
			playerTurn = util.Increment(playerTurn, 0, g.MetaData().NumPlayers-1)
		}
	}

	// We're done, whoever's left with dice is the winner.
	for _, p := range g.Players {
		if g.Board.DiceHidden[p].Count() > 0 {
			place := game.Place{
				Player: p.Player,
				Rank:   1,
			}
			g.AddPlace(place)
			break
		}
	}

	return nil
}

func (g *Game) sendRollEvents() {
	for _, p := range g.Players {
		d := g.Board.DiceHidden[p]
		if d.Count() > 0 {
			e := EventRoll{
				ID:   p.ID,
				Dice: d.Values,
			}
			g.EventLog.Add(e, []game.PlayerID{p.ID})
		}
	}
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
		case EventTypeMove:
			e := EventMove{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		case EventTypeChallenge:
			e := EventChallenge{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		case EventTypeRoll:
			e := EventRoll{}
			json.Unmarshal(event.Data, &e)
			eStr = e
		}

		events[i] = eStr
	}

	return events
}

func (g *Game) reset() {
	g.Game.Reset()
	g.Board = NewBoard(g.Players)
	if g.Comms == nil {
		g.Comms = NewComms(g)
	}
}

func (g *Game) shufflePlayers() {
	util.Shuffle(g.Players)

	for i := 1; i <= 4; i++ {
		g.Players[i-1].Position = i
	}
}

func (g *Game) gameOver() bool {
	count := 0
	for _, p := range g.Players {
		if g.Board.DiceHidden[p].Count() > 0 {
			count++
		}
	}

	return count == 1
}

func (g *Game) setLoser(p *Player) {
	places := []game.Place{}

	for _, player := range g.Players {
		rank := 1
		if player == p {
			// Last place
			rank = 4
		}

		places = append(places, game.Place{
			Player: player.Player,
			Rank:   rank,
		})
	}

	g.SetPlaces(places)
}
