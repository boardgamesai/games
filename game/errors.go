package game

import (
	"errors"
	"fmt"
)

type DQType string

const (
	DQTypeInvalidMove = DQType("badmove")
	DQTypeTimeout     = DQType("timeout")
	DQTypeRuntime     = DQType("runtime")
)

// DQError = DisqualifiedError
type DQError struct {
	ID   PlayerID
	Type DQType
	Msg  string
}

func (e DQError) Error() string {
	return fmt.Sprintf("player %d disqualified (%s): %s", e.ID, e.Type, e.Msg)
}

func (e DQError) Unwrap() error {
	return errors.New(e.Msg)
}
