package util

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Increment handles incrementing a value that wraps around
func Increment(current, min, max int) int {
	if current == max {
		return min
	}

	return current + 1
}

// RandInt returns a random int in the range [min, max]
func RandInt(min, max int) int {
	randInt, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		// Don't bother returning this, if we're here something is deeply wrong
		panic(fmt.Sprintf("error getting random int: %s", err))
	}

	return int(randInt.Int64()) + min
}

// Shuffle random sorts a slice using the Fisher-Yates algorithm.
// It panics if you pass it something other than a slice.
// Code mostly copied from sort.Slice().
func Shuffle[T any](s []T) {
	maxlen := len(s) - 1 // -1 because the last element can only be swapped with itself
	for i := 0; i < maxlen; i++ {
		j := RandInt(i, maxlen)
		s[i], s[j] = s[j], s[i]
	}
}
