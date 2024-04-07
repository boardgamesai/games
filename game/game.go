package game

import (
	"sort"
)

type Game[P PlayerBaseable] struct {
	Name
	Players []P
	EventLog
	output map[PlayerID]string
	places []Place
}

func (g *Game[P]) Reset() {
	g.EventLog.Clear()
	g.output = map[PlayerID]string{}
	g.places = []Place{}
}

func (g *Game[P]) GetPlayers() []*Player {
	players := []*Player{}
	for _, p := range g.Players {
		players = append(players, p.BasePlayer())
	}

	return players
}

func (g *Game[P]) InitPlayers(newfn func() P) {
	g.Players = make([]P, g.MetaData().NumPlayers)
	for i := 0; i < len(g.Players); i++ {
		g.Players[i] = newfn()
	}
}

func (g *Game[P]) MetaData() MetaData {
	return Data[g.Name]
}

func (g *Game[P]) LoggedOutput(id PlayerID) string {
	return g.output[id]
}

func (g *Game[P]) SetOutput(id PlayerID, r Runnable) {
	g.output[id] = r.Stderr()
}

func (g *Game[P]) RawEvents() EventLog {
	return g.EventLog
}

func (g *Game[P]) Places() []Place {
	return g.places
}

func (g *Game[P]) SetPlaces(places []Place) {
	// Let's ensure these are sorted here
	sort.Slice(places, func(i, j int) bool { return places[i].Rank < places[j].Rank })
	g.places = places
}

func (g *Game[P]) AddPlace(place Place) {
	g.places = append(g.places, place)
}

func (g *Game[P]) AddDQErrorID(err *DQError, id PlayerID) *DQError {
	return &DQError{
		ID:   id,
		Type: err.Type,
		Msg:  err.Msg,
	}
}
