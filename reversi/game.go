package reversi

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
	g.Name = game.Reversi

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
			Disc:  player.Disc,
		}
		setupEvent.Players = append(setupEvent.Players, esp)
	}

	g.EventLog.AddNone(setupEvent)

	// Game is over when board is filled or no one has moves left
	playerTurn := 0
	for !g.board.IsFull() {
		player := g.players[playerTurn]
		move, err := g.Comms.GetMove(player)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		flips, err := g.board.ApplyMove(player.Disc, move)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}

		e := EventMove{
			ID:    player.ID,
			Move:  move,
			Flips: flips,
			Score: g.board.Score(),
		}
		g.EventLog.AddAll(e)

		playerTurn = util.Increment(playerTurn, 0, 1)
		if len(g.board.PossibleMoves(g.players[playerTurn].Disc)) == 0 {
			// No moves left for this player, skip them
			playerTurn = util.Increment(playerTurn, 0, 1)
			if len(g.board.PossibleMoves(g.players[playerTurn].Disc)) == 0 {
				// No moves left for this player either! Game is over.
				break
			}
		}
	}

	score := g.board.Score()
	var winner *Player
	if score[Black] > score[White] {
		winner = g.playerDisc(Black)
	} else if score[Black] < score[White] {
		winner = g.playerDisc(White)
	}
	g.setWinner(winner)

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
	g.board = NewBoard()
	if g.Comms == nil {
		g.Comms = NewComms(g)
	}
}

func (g *Game) setWinner(p *Player) {
	var places []game.Place
	scores := g.board.Score()

	if p == nil {
		player1 := g.players[0]
		player2 := g.players[1]
		places = []game.Place{
			{Player: player1.Player, Rank: 1, Tie: true, Score: scores[player1.Disc]},
			{Player: player2.Player, Rank: 1, Tie: true, Score: scores[player2.Disc]},
		}
	} else {
		loser := g.otherPlayer(p)
		places = []game.Place{
			{Player: p.Player, Rank: 1, Tie: false, Score: scores[p.Disc]},
			{Player: loser.Player, Rank: 2, Tie: false, Score: scores[loser.Disc]},
		}
	}

	g.SetPlaces(places)
}

func (g *Game) shufflePlayers() {
	util.Shuffle(g.players)

	discs := []Disc{Black, White}
	for i := 0; i < 2; i++ {
		g.players[i].Order = i + 1
		g.players[i].Disc = discs[i]
	}
}

func (g *Game) otherPlayer(player *Player) *Player {
	if g.players[0] == player {
		return g.players[1]
	}
	return g.players[0]
}

func (g *Game) playerDisc(d Disc) *Player {
	for _, p := range g.players {
		if p.Disc == d {
			return p
		}
	}
	return nil // If we get here, there's problems
}

func (g *Game) String() string {
	return g.board.String()
}
