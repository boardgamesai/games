package fourinarow

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
	game := Game{
		Players: make([]*Player, 2),
		Board:   &Board{},
		Moves:   []MoveLog{},
	}

	return &game
}

func (g *Game) Play() error {
	config, err := game.Config()
	if err != nil {
		return err
	}

	g.ShufflePlayers()

	// Launch the player processes
	for _, player := range g.Players {
		err = player.Run(config.UseSandbox)
		if err != nil {
			return fmt.Errorf("player %s failed to run, err: %s", player, err)
		}
		defer player.Terminate()

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
			g.Winner = g.Players[playerTurn^1]
			return fmt.Errorf("player %s failed to get move, err: %s stderr: %s", player, err, player.Stderr())
		}

		err = g.Board.ApplyMove(player, move)
		if err != nil {
			g.Winner = g.Players[playerTurn^1]
			return fmt.Errorf("player %s committed invalid move: %s err: %s", player, move, err)
		}

		g.Moves = append(g.Moves, MoveLog{Move: move, Order: player.Order})

		if g.Board.HasWinner() {
			g.Winner = player
			break
		}

		// Bitwise XOR with 1 flips 0 -> 1 and 1 -> 0
		playerTurn = playerTurn ^ 1
	}

	return nil
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

func (g *Game) String() string {
	return fmt.Sprintf("%s", g.Board)
}
