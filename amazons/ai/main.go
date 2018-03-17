package main

import "github.com/boardgamesai/games/amazons/ai/driver"

func main() {
	d := driver.New(&AI{})
	d.Run()
}
