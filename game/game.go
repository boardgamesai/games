package game

type Game struct {
	Name
	EventLog
	output map[PlayerID]string
	places []Place
}

func (g *Game) Reset() {
	g.EventLog.Clear()
	g.output = map[PlayerID]string{}
	g.places = []Place{}
}

func (g *Game) MetaData() MetaDataEntry {
	return MetaData[g.Name]
}

func (g *Game) LoggedOutput(id PlayerID) string {
	return g.output[id]
}

func (g *Game) SetOutput(id PlayerID, r Runnable) {
	g.output[id] = r.Stderr()
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
