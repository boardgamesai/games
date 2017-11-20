package game

import (
	"fmt"
	"strings"
)

type Game struct {
	EventLog
	PlayersCached []*Player // We use this to lazily store the game.Players embedded in our game-specific players
	places        []Place
}

func (g *Game) Reset() {
	g.EventLog = EventLog{}
	g.PlayersCached = []*Player{}
	g.places = []Place{}
}

func (g *Game) Places() []Place {
	return g.places
}

func (g *Game) SetPlaces(places []Place) {
	g.places = places
}

func Usage(gameName string, numPlayers int) string {
	players := make([]string, numPlayers)
	for i := 1; i <= numPlayers; i++ {
		players[i-1] = fmt.Sprintf("<player%d>", i)
	}

	return fmt.Sprintf("go run play.go [-n numGames] %s %s", gameName, strings.Join(players, " "))
}
