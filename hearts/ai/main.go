package main

import "github.com/boardgamesai/games/hearts/ai/driver"

func main() {
	d := driver.New(&AI{})
	d.Run()
}
