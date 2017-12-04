package game

type Game struct {
	EventLog
	PlayersCached []*Player // We use this to lazily store the game.Players embedded in our game-specific players
	output        map[int]string
	places        []Place
}

func (g *Game) Reset() {
	g.EventLog.Clear()
	g.PlayersCached = []*Player{}
	g.output = map[int]string{}
	g.places = []Place{}
}

func (g *Game) LoggedOutput(order int) string {
	return g.output[order]
}

func (g *Game) SetOutput(order int, r Runnable) {
	g.output[order] = r.Stderr()
}

func (g *Game) Places() []Place {
	return g.places
}

func (g *Game) SetPlaces(places []Place) {
	g.places = places
}
