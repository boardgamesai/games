package game

import "fmt"

type Playable interface {
	NumPlayers() int
	Play() error
	AddPlayer(name string)
	Players() []*Player
	Events() []fmt.Stringer
	Places() []Place
}
