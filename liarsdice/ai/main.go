package main

import "github.com/boardgamesai/games/liarsdice/ai/driver"

func main() {
	d := driver.New(&AI{})
	d.Run()
}
