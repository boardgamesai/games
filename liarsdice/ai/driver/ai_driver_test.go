package driver

import (
	"encoding/json"
	"testing"

	"github.com/boardgamesai/games/game"
	"github.com/boardgamesai/games/liarsdice"
)

func TestSetup(t *testing.T) {
	players := []*liarsdice.Player{}
	for i := 1; i <= 4; i++ {
		player := liarsdice.Player{
			Player: game.Player{
				ID: game.PlayerID(i),
			},
			Position: i,
		}
		players = append(players, &player)
	}

	m := liarsdice.MessageSetup{
		ID:       2,
		Position: 2,
		Players:  players,
	}
	mJSON, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("error marshaling move: %s", err)
	}

	d := Driver{}
	output, err := d.handleSetup(mJSON)
	if err != nil {
		t.Fatalf("error handling setup: %s", err)
	}

	if string(output) != "\"OK\"" {
		t.Fatalf("got unexpected handleSetup output: %s", output)
	}
}
