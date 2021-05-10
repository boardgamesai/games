package reversi

import "github.com/boardgamesai/games/game"

type MessageSetup struct {
	Disc
	ID       game.PlayerID
	Order    int
	Opponent *Player
}

type MessageMove struct {
	NewEvents []game.Event
}
