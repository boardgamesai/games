package game

import "fmt"

type Playable interface {
	NumPlayers() int
	Play() error
	AddPlayer(name string, r Runnable)
	Players() []*Player
	Events() []fmt.Stringer
	Places() []Place
	LoggedOutput(order int) string
}
