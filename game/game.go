package game

type Game struct {
	EventLog
	output map[int]string
	places []Place
}

func (g *Game) Reset() {
	g.EventLog.Clear()
	g.output = map[int]string{}
	g.places = []Place{}
}

func (g *Game) LoggedOutput(order int) string {
	return g.output[order]
}

func (g *Game) SetOutput(order int, r Runnable) {
	g.output[order] = r.Stderr()
}

func (g *Game) RawEvents() EventLog {
	return g.EventLog
}

func (g *Game) Places() []Place {
	return g.places
}

func (g *Game) SetPlaces(places []Place) {
	g.places = places
}
