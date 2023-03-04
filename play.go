package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/game/factory"
)

func main() {
	numGamesFlag := flag.Int("n", 1, "number of games to play")
	randomFlag := flag.Bool("random", false, "should we play a random game")
	rawEventsFlag := flag.Bool("raw", false, "display raw or formatted event log")
	printBoardFlag := flag.Bool("print", false, "print the board at the end of the game")
	flag.Parse()

	numGames := *numGamesFlag
	playRandom := *randomFlag

	if numGames < 1 {
		log.Fatalf("Invalid number of games: %d\n", numGames)
	}

	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("Usage: %s", usageNoGame())
	}

	gameName := game.Name(args[0])
	g, err := factory.New(gameName)
	if err != nil {
		log.Fatalf("%s", err)
	}

	numPlayers := game.MetaData[gameName].NumPlayers
	filenames := []string{}
	if playRandom {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("%s", err)
		}

		path := fmt.Sprintf("%s/%s/ai/example/random/random.go", pwd, gameName)
		for i := 0; i < numPlayers; i++ {
			filenames = append(filenames, path)
		}
	} else {
		if len(args) != numPlayers+1 { // +1 to accommodate the game name
			log.Fatalf("Usage: %s", usage(gameName, numPlayers))
		}
		filenames = args[1:]
	}

	players := g.Players()
	for i, filename := range filenames {
		players[i].ID = game.PlayerID(i + 1)
		players[i].Name = game.FileNameToPlayerName(filename)
		players[i].Runnable = game.NewRunnablePlayer(string(gameName), filename)
	}

	if numGames == 1 {
		playOneGame(g, gameName, *rawEventsFlag, *printBoardFlag)
	} else {
		playMultipleGames(g, numGames)
	}
}

func playOneGame(g game.Playable, gameName game.Name, showRawEvents, printBoard bool) {
	gameErr := g.Play()

	fmt.Printf("Ordered players:\n")
	for _, player := range g.Players() {
		fmt.Printf("* %s (ID: %d)\n", player.Name, player.ID)
	}

	fmt.Println()

	if showRawEvents {
		for i, event := range g.RawEvents() {
			fmt.Printf("%d. %s\n", i+1, event)
		}
	} else {
		for i, event := range g.Events() {
			fmt.Printf("%d. %s\n", i+1, event)
		}
	}

	fmt.Println("\nFinish places:")

	for _, place := range g.Places() {
		tie := ""
		if place.Tie {
			tie = " (tie)"
		}
		fmt.Printf("%d.%s %s (ID: %d)", place.Rank, tie, place.Player.Name, place.Player.ID)
		if game.MetaData[gameName].HasScore {
			fmt.Printf(": %d", place.Score)
		}
		fmt.Println()
	}

	if gameErr != nil {
		fmt.Printf("*** game ended with error: %s\n", gameErr)
	}

	if printBoard {
		fmt.Printf("\n%s\n", g)
	}

	printLoggedOutput(g)
}

func playMultipleGames(g game.Playable, numGames int) {
	// Grab a copy of the game's players in the original order. The game will shuffle them,
	// but we want to know the original order for reporting purposes.
	players := make([]*game.Player, len(g.Players()))
	copy(players, g.Players())

	outcomes := map[game.PlayerID]map[int]int{}
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

func usage(gameName game.Name, numPlayers int) string {
	players := make([]string, numPlayers)
	for i := 1; i <= numPlayers; i++ {
		players[i-1] = fmt.Sprintf("<player%d>", i)
	}

	return fmt.Sprintf("go run play.go [-n numGames] %s %s", gameName, strings.Join(players, " "))
}

func usageNoGame() string {
	return "go run play.go [-n numGames] <game> <player1> <player2> ..."
}

func printLoggedOutput(g game.Playable) {
	for _, player := range g.Players() {
		loggedOutput := g.LoggedOutput(player.ID)
		if loggedOutput != "" {
			fmt.Printf("Player %d logged output:\n", player.ID)
			fmt.Printf("%s\n", loggedOutput)
		}
	}
}

func printSummaryTotals(players []*game.Player, outcomes map[game.PlayerID]map[int]int) {
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
