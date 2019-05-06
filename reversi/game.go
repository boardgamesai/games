package reversi

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

	// Game is over when board is filled or no one has moves left
	playerTurn := 0
	for !g.board.IsFull() {
		player := g.players[playerTurn]
		move, err := g.Comms.GetMove(player)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.board.ApplyMove(player.Disc, move)
		if err != nil {
			g.setWinner(g.otherPlayer(player))
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}

		e := EventMove{
			Order: player.Order,
			Disc:  player.Disc,
			Move:  move,
		}
		g.EventLog.Add(e, game.AllPlayers)

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

func (g *Game) AddPlayer(filepath string, r game.Runnable) {
	if r == nil {
		r = game.NewRunnablePlayer("reversi", filepath)
	}

	player := Player{
		Player: game.Player{
			ID:   uuid.NewRandom().String(), // HACK - TODO use actual IDs
			Name: game.FileNameToPlayerName(filepath),
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
	var places []game.Place
	scores := g.board.Score()

	if p == nil {
		player1 := g.players[0]
		player2 := g.players[1]
		places = []game.Place{
			{Player: player1.Player, Rank: 1, Tie: true, Score: scores[player1.Disc], HasScore: true},
			{Player: player2.Player, Rank: 1, Tie: true, Score: scores[player2.Disc], HasScore: true},
		}
	} else {
		loser := g.otherPlayer(p)
		places = []game.Place{
			{Player: p.Player, Rank: 1, Tie: false, Score: scores[p.Disc], HasScore: true},
			{Player: loser.Player, Rank: 2, Tie: false, Score: scores[loser.Disc], HasScore: true},
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
	return fmt.Sprintf("%s", g.board)
}
