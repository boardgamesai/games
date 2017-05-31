package tictactoe

import (
	"fmt"

	"github.com/boardgamesai/games/util"
)

type Game struct {
	Players []*Player
	Board   *Board
	Moves   []MoveLog
	Turn    int
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

	// Launch the player processes
	for _, player := range g.Players {
		err := player.Run(config.UseSandbox)
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}
	}

	// Game is over when someone wins or board is filled
	for !g.Board.IsFull() {
		player := g.CurrentPlayer()
		move, err := player.GetMove(g.Board)
		if err != nil {
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.ApplyMove(move, player)
		if err != nil {
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}

		g.Moves = append(g.Moves, MoveLog{move, player})

		if g.Board.HasWinner() {
			g.Winner = player
			break
		}

		g.Turn = g.nextTurn()
	}

	return nil
}

func (g *Game) ApplyMove(m Move, p *Player) error {
	return g.Board.Set(m.Col, m.Row, p.Symbol)
}

func (g *Game) AddPlayer(symbol string, playerPath string, aiPath string) {
	player := Player{
		Symbol:     symbol,
		PlayerPath: playerPath,
		AIPath:     aiPath,
	}
	g.Players = append(g.Players, &player)
}

func (g *Game) CurrentPlayer() *Player {
	return g.Players[g.Turn]
}

func (g *Game) String() string {
	return fmt.Sprintf("%s", g.Board)
}

func (g *Game) nextTurn() int {
	if g.Turn == 0 {
		return 1
	}
	return 0
}
