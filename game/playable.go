package game

import "fmt"

type Playable interface {
	Play() error
	GetPlayers() []*Player
	Events() []fmt.Stringer
	RawEvents() EventLog
	Places() []Place
	LoggedOutput(id PlayerID) string
}
