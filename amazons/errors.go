package amazons

import "fmt"

type OutOfBoundsError struct {
	Move Move
}

func (e OutOfBoundsError) Error() string {
	return fmt.Sprintf("out of bounds: %s", e.Move)
}

type InvalidFromError struct {
	Move Move
}

func (e InvalidFromError) Error() string {
	return fmt.Sprintf("invalid from: %s", e.Move)
}

type InvalidToError struct {
	Move Move
}

func (e InvalidToError) Error() string {
	return fmt.Sprintf("invalid to: %s", e.Move)
}

type InvalidArrowError struct {
	Move Move
}

func (e InvalidArrowError) Error() string {
	return fmt.Sprintf("invalid arrow: %s", e.Move)
}
