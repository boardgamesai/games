package game

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
