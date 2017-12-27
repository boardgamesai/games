package main

import "github.com/boardgamesai/games/reversi/ai/driver"

func main() {
	d := driver.New(&AI{})
	d.Run()
}
