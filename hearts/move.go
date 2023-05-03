package hearts

import "github.com/boardgamesai/games/game/elements/card"

type PassMove struct {
	Cards []card.Card // Must contain exactly three cards
}

type PlayMove struct {
	Card card.Card
}
