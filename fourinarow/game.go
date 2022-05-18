package fourinarow

import (
	"encoding/json"
	"fmt"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/util"
)

type Game struct {
	game.Game
	Comms   AIComms
	players []*Player
	board   *Board
}

func New() *Game {
	g := Game{
		players: []*Player{},
	}
	g.Name = game.FourInARow

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
		}
		setupEvent.Players = append(setupEvent.Players, esp)
	}

	g.EventLog.AddNone(setupEvent)

	// Game is over when someone wins or board is filled
	playerTurn := 0
	for !g.board.IsFull() {
		player := g.players[playerTurn]
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

		row, err := g.board.ApplyMove(player.Order, move)
		e := EventMove{
			ID:   player.ID,
			Move: move,
			Row:  row,
		}
		hasWinner, winCoords := g.board.HasWinner()
		if hasWinner {
			e.WinCoords = winCoords
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
	g.board = &Board{}
	if g.Comms == nil {
		g.Comms = NewComms(g)
	}
}

func (g *Game) setWinner(p *Player) {
	var places []game.Place

	if p == nil {
		places = []game.Place{
			{Player: g.players[0].Player, Rank: 1, Tie: true},
			{Player: g.players[1].Player, Rank: 1, Tie: true},
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
	util.Shuffle(g.players)

	for i := 0; i < 2; i++ {
		g.players[i].Order = i + 1
	}
}

func (g *Game) otherPlayer(player *Player) *Player {
	if g.players[0] == player {
		return g.players[1]
	}
	return g.players[0]
}

func (g *Game) String() string {
	return g.board.String()
}
