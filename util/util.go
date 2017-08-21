package util

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"os"
)

// RandInt returns a random int in the range [min, max]
func RandInt(min, max int) int {
	randInt, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		// Don't bother returning this, if we're here something is deeply wrong
		panic(fmt.Sprintf("error getting random int: %s", err))
	}

	return int(randInt.Int64()) + min
}

// CoinFlip is syntactic sugar for picking a random 0 or 1
func CoinFlip() bool {
	return RandInt(0, 1) == 1
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
