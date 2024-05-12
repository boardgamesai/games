package tictactoe

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
	g.Name = game.TicTacToe
	g.InitPlayers(NewPlayer)
	return &g
}

func (g *Game) Play() error {
	// Wipe out any previous state
	g.reset()

	// Decide who is X and goes first
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
			ID:     player.ID,
			Order:  player.Order,
			Symbol: player.Symbol,
		}
		setupEvent.Players = append(setupEvent.Players, esp)
	}

	g.EventLog.AddNone(setupEvent)

	// Game is over when someone wins or board is filled
	playerTurn := 0
	for !g.Board.IsFull() {
		player := g.Players[playerTurn]
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

		err = g.Board.ApplyMove(player.Symbol, move)

		// Regardless of whether the move was valid or not, we add it to the log
		e := EventMove{
			ID:   player.ID,
			Move: move,
		}
		hasWinner, winMoves := g.Board.HasWinner()
		if hasWinner {
			e.WinMoves = winMoves
		}
		g.EventLog.AddAll(e)

		// Now see if the move was valid
		if err != nil {
			// Disqualification - game over and the other player wins
			g.setWinner(g.otherPlayer(player))
			return game.DQError{
				ID:   player.ID,
				Type: game.DQTypeInvalidMove,
				Msg:  err.Error(),
			}
		}

		// If there was a winner, we stop
		if hasWinner {
			g.setWinner(player)
			break
		}

		playerTurn = util.Increment(playerTurn, 0, 1)
	}

	if len(g.Places()) == 0 {
		// No winner, so this is a tie.
		g.setWinner(nil)
	}

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
	g.Board = &Board{}
	if g.Comms == nil {
		g.Comms = NewComms(g)
	}
}

func (g *Game) setWinner(p *Player) {
	var places []game.Place

	if p == nil {
		places = []game.Place{
			{Player: g.Players[0].Player, Rank: 1, Tie: true},
			{Player: g.Players[1].Player, Rank: 1, Tie: true},
		}
	} else {
		places = []game.Place{
			{Player: p.Player, Rank: 1, Tie: false},
			{Player: g.otherPlayer(p).Player, Rank: 2, Tie: false},
		}
	}

	g.SetPlaces(places)
}

func (g *Game) shufflePlayers() {
	util.Shuffle(g.Players)

	symbols := []string{"X", "O"}
	for i := 0; i < 2; i++ {
		g.Players[i].Order = i + 1
		g.Players[i].Symbol = symbols[i]
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
