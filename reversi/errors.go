package reversi

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

type IllegalMoveError struct {
	Move Move
}

func (e IllegalMoveError) Error() string {
	return fmt.Sprintf("illegal move: %s", e.Move)
}
