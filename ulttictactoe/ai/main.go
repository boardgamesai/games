package main

import "github.com/boardgamesai/games/ulttictactoe/ai/driver"

func main() {
	d := driver.New(&AI{})
	d.Run()
}
