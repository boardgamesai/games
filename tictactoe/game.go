package tictactoe

import (
	"fmt"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/util"
)

type Game struct {
	Players []*Player
	Board   *Board
	Moves   []MoveLog
	Winner  *Player
}

func NewGame() *Game {
	g := Game{
		Players: []*Player{},
	}
	g.Reset()
	return &g
}

func (g *Game) Reset() {
	g.Board = &Board{}
	g.Moves = []MoveLog{}
	g.Winner = nil
}

func (g *Game) Play() error {
	config, err := game.Config()
	if err != nil {
		return err
	}

	// Wipe out any previous state
	g.Reset()

	// Decide who is X and goes first
	g.ShufflePlayers()

	// Launch the player processes
	for _, player := range g.Players {
		// This copies files to a tmp dir, runs it, and sends a heartbeat message to verify.
		err = player.Run(config)
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}
		defer player.CleanUp()

		// This initializes the game state for this player.
		err = player.Setup(g)
		if err != nil {
			return fmt.Errorf("player %s failed to setup, err: %s", player, err)
		}
	}

	// Game is over when someone wins or board is filled
	playerTurn := 0
	for !g.Board.IsFull() {
		player := g.Players[playerTurn]
		move, err := player.GetMove(g)
		if err != nil {
			g.Winner = g.OtherPlayer(player)
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.Board.IsValidMove(move)
		if err != nil {
			g.Winner = g.OtherPlayer(player)
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}

		g.Board.Grid[move.Col][move.Row] = player.Symbol
		g.Moves = append(g.Moves, MoveLog{Move: move, Order: player.Order})

		if g.Board.HasWinner() {
			g.Winner = player
			break
		}

		playerTurn = util.Increment(playerTurn, 0, 1)
	}

	return nil
}

func (g *Game) AddPlayer(name string) {
	basePlayer := game.NewPlayer("tictactoe", name)
	player := Player{
		Name:   name,
		Player: *basePlayer,
	}
	g.Players = append(g.Players, &player)
}

func (g *Game) ShufflePlayers() {
	util.Shuffle(g.Players)

	symbols := []string{"X", "O"}
	for i := 0; i < 2; i++ {
		g.Players[i].Order = i + 1
		g.Players[i].Symbol = symbols[i]
	}
}

// GetNewMovesForPlayer crawls the move log and finds all the moves since p's last move
func (g *Game) GetNewMovesForPlayer(p *Player) []MoveLog {
	moves := []MoveLog{}

	for i := len(g.Moves) - 1; i >= 0; i-- {
		if g.Moves[i].Order == p.Order {
			break
		}
		moves = append([]MoveLog{g.Moves[i]}, moves...)
	}

	return moves
}

func (g *Game) OtherPlayer(player *Player) *Player {
	if g.Players[0] == player {
		return g.Players[1]
	}
	return g.Players[0]
}

func (g *Game) String() string {
	return fmt.Sprintf("%s", g.Board)
}
