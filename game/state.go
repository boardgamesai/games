package game

type State struct {
	ID        PlayerID
	AllEvents []Event
	NewEvents []Event
}

func (s *State) AddEvents(events []Event) {
	s.AllEvents = append(s.AllEvents, events...)
	s.NewEvents = events
}
