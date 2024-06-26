package amazons

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
	g.Name = game.Amazons
	g.InitPlayers(NewPlayer)
	return &g
}

func (g *Game) Play() error {
	// Wipe out any previous state
	g.reset()

	// Decide who goes first
	g.shufflePlayers()

	// We need to write down our setup
	setupEvent := EventSetup{
		Players: []EventSetupPlayer{},
	}

	// Launch the player processes
	for _, player := range g.Players {
		defer player.CleanUp()
		defer g.SetOutput(player.ID, player)

		// This copies files to a tmp dir, runs it, and sends a heartbeat message to verify.
		err := player.Run()
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}

		// This initializes the game state for this player.
		err = g.Comms.Setup(player, g.otherPlayer(player))
		if err != nil {
			return fmt.Errorf("player %s failed to setup, err: %s", player, err)
		}

		// Keep track of our setup
		esp := EventSetupPlayer{
			ID:    player.ID,
			Order: player.Order,
			Color: player.Color,
		}
		setupEvent.Players = append(setupEvent.Players, esp)
	}

	g.EventLog.AddNone(setupEvent)

	// Game is over when someone can't move - a draw is impossible
	playerTurn := 0
	for {
		player := g.Players[playerTurn]
		if !g.Board.CanMove(player.Color) {
			break // Someone can't move, game over
		}

		move, err := g.Comms.GetMove(player)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
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

		err = g.Board.ApplyMove(player.Color, move)

		e := EventMove{
			ID:   player.ID,
			Move: move,
		}
		g.EventLog.AddAll(e)

		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return game.DQError{
				ID:   player.ID,
				Type: game.DQTypeInvalidMove,
				Msg:  err.Error(),
			}
		}

		playerTurn = util.Increment(playerTurn, 0, 1)
	}

	// Whoever's turn it is when we get here is the loser
	g.setWinner(g.otherPlayer(g.Players[playerTurn]))

	return nil
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
		}

		events[i] = eStr
	}

	return events
}

func (g *Game) reset() {
	g.Game.Reset()
	g.Board = NewBoard()
	if g.Comms == nil {
		g.Comms = NewComms(g)
	}
}

func (g *Game) setWinner(p *Player) {
	places := []game.Place{
		{Player: p.Player, Rank: 1},
		{Player: g.otherPlayer(p).Player, Rank: 2},
	}
	g.SetPlaces(places)
}

func (g *Game) shufflePlayers() {
	util.Shuffle(g.Players)

	colors := []SpaceType{White, Black}
	for i := 0; i < 2; i++ {
		g.Players[i].Order = i + 1
		g.Players[i].Color = colors[i]
	}
}

func (g *Game) otherPlayer(player *Player) *Player {
	if g.Players[0] == player {
		return g.Players[1]
	}
	return g.Players[0]
}

func (g *Game) String() string {
	return g.Board.String()
}
