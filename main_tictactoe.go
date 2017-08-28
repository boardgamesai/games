package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/boardgamesai/games/tictactoe"
)

func main() {
	numGamesFlag := flag.Int("n", 1, "number of games to play")
	flag.Parse()

	if len(flag.Args()) != 2 {
		log.Fatalf("Usage: go run main_tictactoe.go [-n <num games>] <player1> <player2>\n")
	}

	numGames := *numGamesFlag
	if numGames < 1 {
		log.Fatalf("Invalid number of games: %d\n", numGames)
	}

	g := tictactoe.NewGame()
	for _, playerName := range flag.Args() {
		g.AddPlayer(playerName)
	}

	// Grab a copy of the game's players in the original order. The game will shuffle them,
	// but we want to know the original order for reporting purposes.
	players := make([]*tictactoe.Player, len(g.Players))
	copy(players, g.Players)

	if numGames == 1 {
		err := g.Play()
		if err != nil {
			fmt.Printf("game ended with error: %s\n", err)
			return
		}

		for _, player := range g.Players {
			fmt.Printf("Player %d: %s\n", player.Order, player.Name)
		}

		fmt.Println()

		for i, log := range g.Moves {
			fmt.Printf("%d. %s\n", i+1, log)
		}

		fmt.Printf("\n%s\n", g)

		if g.Winner != nil {
			fmt.Printf("%s wins!\n", g.Winner)
		} else {
			fmt.Printf("Game is a draw.\n")
		}
		fmt.Println()

		for _, player := range g.Players {
			loggedOutput := player.Stderr()
			if loggedOutput != "" {
				fmt.Printf("Player %d logged output:\n", player.Order)
				fmt.Printf("%s\n", loggedOutput)
			}
		}
	} else {
		outcomes := map[*tictactoe.Player]map[string]int{}
		for _, player := range players {
			outcomes[player] = map[string]int{"win": 0, "lose": 0, "draw": 0}
		}

		for i := 1; i <= numGames; i++ {
			fmt.Printf("playing game %d...\n", i)

			err := g.Play()
			if err != nil {
				fmt.Printf("game %d ended with error: %s\n", i, err)
				continue
			}

			if g.Winner != nil {
				// Somebody won.
				outcomes[g.Winner]["win"]++
				for _, player := range g.Players {
					// Find the non-winner aka the loser
					if player != g.Winner {
						outcomes[player]["lose"]++
						break
					}
				}
			} else {
				// This is a draw.
				for _, player := range g.Players {
					outcomes[player]["draw"]++
				}
			}
		}

		for _, player := range players {
			fmt.Printf("%s: win: %d lose: %d draw: %d\n", player.Name, outcomes[player]["win"],
				outcomes[player]["lose"], outcomes[player]["draw"])
		}
	}
}
