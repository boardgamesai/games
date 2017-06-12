package tictactoe

import (
	"fmt"

	"github.com/boardgamesai/games/util"
)

type Players []*Player

func (p Players) Shuffle() {
	if util.CoinFlip() {
		temp := p[0]
		p[0] = p[1]
		p[1] = temp
	}
}

type Game struct {
	Players
	Board  *Board
	Moves  []MoveLog
	Winner *Player
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
	g.Players.Shuffle()
	g.Players[0].Symbol = "X"
	g.Players[1].Symbol = "O"

	// Launch the player processes
	for _, player := range g.Players {
		err := player.Run(config.UseSandbox)
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}
	}

	// Game is over when someone wins or board is filled
	player := g.Players[0]
	for !g.Board.IsFull() {
		move, err := player.GetMove(g.Board)
		if err != nil {
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.Board.IsValidMove(move)
		if err != nil {
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}
		g.Board.Grid[move.Col][move.Row] = player.Symbol

		g.Moves = append(g.Moves, MoveLog{move, player})

		if g.Board.HasWinner() {
			g.Winner = player
			break
		}

		if player == g.Players[0] {
			player = g.Players[1]
		} else {
			player = g.Players[0]
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

func (g *Game) String() string {
	return fmt.Sprintf("%s", g.Board)
}
