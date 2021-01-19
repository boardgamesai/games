package tictactoe

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
	g.Name = game.TicTacToe

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

	// Decide who is X and goes first
	g.shufflePlayers()

	// Launch the player processes
	for _, player := range g.players {
		defer player.CleanUp()
		defer g.SetOutput(player.Order, player)

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
	}

	// Game is over when someone wins or board is filled
	playerTurn := 0
	for !g.board.IsFull() {
		player := g.players[playerTurn]
		move, err := g.Comms.GetMove(player)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.board.ApplyMove(player.Symbol, move)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return game.DQError{
				ID:  player.ID,
				Err: err,
			}
		}

		e := EventMove{
			Order:  player.Order,
			Symbol: player.Symbol,
			Move:   move,
		}
		g.EventLog.Add(e, game.AllPlayers)

		if g.board.HasWinner() {
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

	symbols := []string{"X", "O"}
	for i := 0; i < 2; i++ {
		g.players[i].Order = i + 1
		g.players[i].Symbol = symbols[i]
	}
}

func (g *Game) otherPlayer(player *Player) *Player {
	if g.players[0] == player {
		return g.players[1]
	}
	return g.players[0]
}

func (g *Game) String() string {
	return fmt.Sprintf("%s", g.board)
}
