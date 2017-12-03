package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/game/factory"
)

func main() {
	numGamesFlag := flag.Int("n", 1, "number of games to play")
	flag.Parse()

	numGames := *numGamesFlag
	if numGames < 1 {
		log.Fatalf("Invalid number of games: %d\n", numGames)
	}

	args := flag.Args()
	gameName := args[0]
	g, err := factory.New(gameName)
	if err != nil {
		log.Fatalf("%s", err)
	}

	if len(args) != (g.NumPlayers() + 1) { // +1 to accomodate the game name
		log.Fatalf("Usage: %s", usage(gameName, g.NumPlayers()))
	}

	for _, playerName := range args[1:] {
		g.AddPlayer(playerName)
	}

	if numGames == 1 {
		err := g.Play()
		if err != nil {
			fmt.Printf("game ended with error: %s\n", err)
			return
		}

		for _, player := range g.Players() {
			fmt.Printf("Player %d: %s\n", player.Order, player.Name)
		}

		fmt.Println()

		for i, event := range g.Events() {
			fmt.Printf("%d. %s\n", i+1, event)
		}

		fmt.Println("\nFinish places:")

		for _, place := range g.Places() {
			tie := ""
			if place.Tie {
				tie = " (tie)"
			}
			fmt.Printf("%d.%s %s (%d)\n", place.Rank, tie, place.Player.Name, place.Player.Order)
		}

		for _, player := range g.Players() {
			loggedOutput := player.Stderr()
			if loggedOutput != "" {
				fmt.Printf("Player %d logged output:\n", player.Order)
				fmt.Printf("%s\n", loggedOutput)
			}
		}
	} else {
		// Grab a copy of the game's players in the original order. The game will shuffle them,
		// but we want to know the original order for reporting purposes.
		players := make([]*game.Player, len(g.Players()))
		copy(players, g.Players())

		outcomes := map[string]map[int]int{}
		for _, player := range players {
			outcomes[player.ID] = map[int]int{}
		}

		for i := 1; i <= numGames; i++ {
			fmt.Printf("playing game %d...\n", i)

			err := g.Play()
			if err != nil {
				fmt.Printf("game %d ended with error: %s\n", i, err)
				continue
			}

			for _, place := range g.Places() {
				outcomes[place.Player.ID][place.Rank]++
			}
		}

		fmt.Println()
		printSummaryTotals(players, outcomes)
	}
}

func usage(gameName string, numPlayers int) string {
	players := make([]string, numPlayers)
	for i := 1; i <= numPlayers; i++ {
		players[i-1] = fmt.Sprintf("<player%d>", i)
	}

	return fmt.Sprintf("go run play.go [-n numGames] %s %s", gameName, strings.Join(players, " "))
}

func printSummaryTotals(players []*game.Player, outcomes map[string]map[int]int) {
	// Let's find the longest player name, so we can pad them appropriately
	maxlen := 0
	for _, player := range players {
		if len(player.Name) > maxlen {
			maxlen = len(player.Name)
		}
	}
	nameFormat := fmt.Sprintf("%%%ds", maxlen)

	// Let's also make a map of max total widths per rank, for similar padding reasons
	rankWidths := map[int]string{}
	for i := 1; i <= len(players); i++ {
		maxlen = 0
		for _, rankMap := range outcomes {
			// Convert score to string, so we can get its length
			scoreStr := fmt.Sprintf("%d", rankMap[i])
			if len(scoreStr) > maxlen {
				maxlen = len(scoreStr)
			}
		}
		rankWidths[i] = fmt.Sprintf("%%%dd", maxlen)
	}

	fmt.Println("Finish places summary:")

	for _, player := range players {
		totals := []string{}
		for i := 1; i <= len(players); i++ {
			totals = append(totals, fmt.Sprintf("%d:"+rankWidths[i], i, outcomes[player.ID][i]))
		}

		fmt.Printf(nameFormat+": %s\n", player.Name, strings.Join(totals, ", "))
	}
}
