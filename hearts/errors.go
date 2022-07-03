package hearts

import "fmt"

type InvalidPassError struct {
	Move PassMove
	Msg  string
}

func (e InvalidPassError) Error() string {
	return fmt.Sprintf("invalid pass: %s", e.Move)
}

// Played a card not in their hand
type InvalidPlayError struct {
	Move PlayMove
}

func (e InvalidPlayError) Error() string {
	return fmt.Sprintf("invalid play: %s", e.Move)
}

// Played a card that the rules prohibit
type IllegalPlayError struct {
	Move PlayMove
}

func (e IllegalPlayError) Error() string {
	return fmt.Sprintf("illegal play: %s", e.Move)
}
