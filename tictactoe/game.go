package tictactoe

import (
	"fmt"

	"github.com/boardgamesai/games/util"
)

type Game struct {
	Players []*Player
	Board   *Board
	Moves   []MoveLog
	Winner  *Player
}

func NewGame() *Game {
	game := Game{
		Players: []*Player{},
		Board:   &Board{},
		Moves:   []MoveLog{},
	}

	return &game
}

func (g *Game) Play() error {
	config, err := util.Config()
	if err != nil {
		return err
	}

	// Decide who is X and goes first
	g.ShufflePlayers()

	// Launch the player processes
	for _, player := range g.Players {
		err := player.Run(config.UseSandbox)
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}
	}

	// Game is over when someone wins or board is filled
	playerTurn := 0
	for !g.Board.IsFull() {
		player := g.Players[playerTurn]
		move, err := player.GetMove(g.Board)
		if err != nil {
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.Board.IsValidMove(move)
		if err != nil {
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}
		g.Board.Grid[move.Col][move.Row] = player.Symbol

		g.Moves = append(g.Moves, MoveLog{Move: move, Order: playerTurn + 1})

		if g.Board.HasWinner() {
			g.Winner = player
			break
		}

		if playerTurn == 0 {
			playerTurn = 1
		} else {
			playerTurn = 0
		}
	}

	return nil
}

func (g *Game) AddPlayer(name string, playerPath string, aiPath string) {
	player := Player{
		Name:       name,
		PlayerPath: playerPath,
		AIPath:     aiPath,
	}
	g.Players = append(g.Players, &player)
}

func (g *Game) ShufflePlayers() {
	if util.CoinFlip() {
		temp := g.Players[0]
		g.Players[0] = g.Players[1]
		g.Players[1] = temp
	}
	for i := 1; i <= 2; i++ {
		g.Players[i-1].Order = i
	}

	g.Players[0].Symbol = "X"
	g.Players[1].Symbol = "O"
}

func (g *Game) String() string {
	return fmt.Sprintf("%s", g.Board)
}
