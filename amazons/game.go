package amazons

import (
	"encoding/json"
	"fmt"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/util"
	"github.com/pborman/uuid"
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
	g.reset()
	return &g
}

func (g *Game) NumPlayers() int {
	return 2
}

func (g *Game) Play() error {
	// Wipe out any previous state
	g.reset()

	// Decide who goes first
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

	// Game is over when someone can't move - a draw is impossible
	playerTurn := 0
	for {
		player := g.players[playerTurn]
		if !g.board.CanMove(player.Color) {
			break // Someone can't move, game over
		}

		move, err := g.Comms.GetMove(player)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.board.ApplyMove(player.Color, move)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}

		e := EventMove{
			Order: player.Order,
			Color: player.Color,
			Move:  move,
		}
		g.EventLog.Add(e, game.AllPlayers)

		playerTurn = util.Increment(playerTurn, 0, 1)
	}

	// Whoever's turn it is when we get here is the loser
	g.setWinner(g.otherPlayer(g.players[playerTurn]))

	return nil
}

func (g *Game) AddPlayer(name string, r game.Runnable) {
	if r == nil {
		r = game.NewRunnablePlayer("amazons", name)
	}

	player := Player{
		Player: game.Player{
			ID:   uuid.NewRandom().String(), // HACK - TODO use actual IDs
			Name: name,
		},
		Runnable: r,
	}
	g.players = append(g.players, &player)
}

func (g *Game) Players() []*game.Player {
	if len(g.PlayersCached) == 0 {
		for _, player := range g.players {
			g.PlayersCached = append(g.PlayersCached, &(player.Player))
		}
	}
	return g.PlayersCached
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
	g.board = NewBoard()
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
	util.Shuffle(g.players)

	colors := []SpaceType{White, Black}
	for i := 0; i < 2; i++ {
		g.players[i].Order = i + 1
		g.players[i].Color = colors[i]
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