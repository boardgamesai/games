package game

import "fmt"

type Playable interface {
	Play() error
	Players() []*Player
	Events() []fmt.Stringer
	RawEvents() EventLog
	Places() []Place
	LoggedOutput(id PlayerID) string
}
