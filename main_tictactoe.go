package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/boardgamesai/games/game"
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

	config, err := game.Config()
	if err != nil {
		log.Fatalf("Could not read config: %s err: %s", game.ConfigPath, err)
	}

	players := []*tictactoe.Player{}
	for _, playerName := range flag.Args() {
		runnablePlayer, err := game.NewRunnablePlayer(config, "tictactoe", playerName)
		if err != nil {
			log.Fatalf("Error setting up player %s: %s", playerName, err)
		}
		defer os.RemoveAll(filepath.Dir(runnablePlayer.PlayerPath))

		player := tictactoe.Player{
			Name:           playerName,
			RunnablePlayer: *runnablePlayer,
		}
		players = append(players, &player)
	}

	if numGames == 1 {
		game := tictactoe.NewGame()
		game.Players = players

		err = game.Play()
		if err != nil {
			fmt.Printf("game ended with error: %s\n", err)
		}

		for _, player := range game.Players {
			fmt.Printf("Player %d: %s\n", player.Order, player.Name)
		}

		fmt.Println()

		for i, log := range game.Moves {
			fmt.Printf("%d. %s\n", i+1, log)
		}

		fmt.Printf("\n%s\n", game)

		if game.Winner != nil {
			fmt.Printf("%s wins!\n", game.Winner)
		} else {
			fmt.Printf("Game is a draw.\n")
		}
		fmt.Println()

		for _, player := range game.Players {
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
			game := tictactoe.NewGame()
			copy(game.Players, players) // Copy so we preserve original command line order for end tally

			err = game.Play()
			if err != nil {
				fmt.Printf("game %d ended with error: %s\n", i, err)
			}

			if game.Winner != nil {
				// Somebody won.
				outcomes[game.Winner]["win"]++
				for _, player := range game.Players {
					// Find the non-winner aka the loser
					if player != game.Winner {
						outcomes[player]["lose"]++
						break
					}
				}
			} else {
				// This is a draw.
				for _, player := range game.Players {
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
