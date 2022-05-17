package fourinarow

import "fmt"

type OutOfBoundsError struct {
	Move Move
}

func (e OutOfBoundsError) Error() string {
	return fmt.Sprintf("out of bounds: %s", e.Move)
}

type ColumnFullError struct {
	Move Move
}

func (e ColumnFullError) Error() string {
	return fmt.Sprintf("not empty: %s", e.Move)
}
