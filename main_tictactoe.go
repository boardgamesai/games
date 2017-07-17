package main

import (
	"fmt"
	"log"
	"os"

	"github.com/boardgamesai/games/tictactoe"
	"github.com/boardgamesai/games/util"
	"github.com/pborman/uuid"
)

const playerFile = "player_tictactoe.go"

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: go run main_tictactoe.go <player1> <player2>\n")
	}

	config, err := util.Config()
	if err != nil {
		log.Fatalf("Could not read config: %s err: %s", util.ConfigPath, err)
	}

	game := tictactoe.NewGame()

	for _, player := range os.Args[1:] {
		aiSrcPath := os.Getenv("GOPATH") + config.PlayerDir + "/tictactoe/" + player + "/" + player + ".go"
		if _, err := os.Stat(aiSrcPath); os.IsNotExist(err) {
			log.Fatalf("Player file does not exist: %s", aiSrcPath)
		}

		// First create the tmp dir for the player
		tmpDir := os.Getenv("GOPATH") + config.TmpDir + "/" + uuid.NewRandom().String()
		err = os.Mkdir(tmpDir, 0700)
		if err != nil {
			log.Fatalf("Could not create tmp dir: %s for player: %s err: %s", tmpDir, player, err)
		}
		defer os.RemoveAll(tmpDir)

		// Next copy over the base player file
		playerDestPath := tmpDir + "/" + playerFile
		err = util.CopyFile(playerFile, playerDestPath)
		if err != nil {
			log.Fatalf("Could not copy %s to %s", playerFile, playerDestPath)
		}

		// Now copy over the AI-specific file
		aiDestPath := tmpDir + "/" + player + ".go"
		err = util.CopyFile(aiSrcPath, aiDestPath)
		if err != nil {
			log.Fatalf("Could not copy %s to %s", aiSrcPath, aiDestPath)
		}

		game.Players = append(game.Players, &tictactoe.Player{
			Name:       player,
			PlayerPath: playerDestPath,
			AIPath:     aiDestPath,
		})
	}

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
}
