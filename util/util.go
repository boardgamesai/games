package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
	"reflect"
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
func Shuffle(slice interface{}) {
	rv := reflect.ValueOf(slice)
	maxlen := rv.Len() - 1 // -1 because the last element can only be swapped with itself
	swap := reflect.Swapper(slice)

	for i := 0; i < maxlen; i++ {
		swap(i, RandInt(i, maxlen))
	}
}

// CopyFile copies a file from source to dest
func CopyFile(srcPath string, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
