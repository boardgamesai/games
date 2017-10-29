package main

import "github.com/boardgamesai/games/tictactoe/ai/driver"

func main() {
	d := driver.New(&AI{})
	d.Run()
}
