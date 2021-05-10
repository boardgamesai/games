package game

import "fmt"

// DQError = DisqualifiedError
type DQError struct {
	ID  PlayerID
	Err error
}

func (e DQError) Error() string {
	return fmt.Sprintf("player %d disqualified: %s", e.ID, e.Err.Error())
}

func (e DQError) Unwrap() error {
	return e.Err
}
