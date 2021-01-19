package tictactoe

import "fmt"

type OutOfBoundsError struct {
	Move Move
}

func (e OutOfBoundsError) Error() string {
	return fmt.Sprintf("out of bounds: %s", e.Move)
}

type NotEmptyError struct {
	Move Move
}

func (e NotEmptyError) Error() string {
	return fmt.Sprintf("not empty: %s", e.Move)
}
