package main

import "github.com/boardgamesai/games/fourinarow/ai/driver"

func main() {
	d := driver.New(&AI{})
	d.Run()
}
