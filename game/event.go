package game

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type string
	Data json.RawMessage
	Show map[PlayerID]bool `json:"-"` // Player IDs who should be shown this event
	Seen map[PlayerID]bool `json:"-"` // Player IDs who have seen this event
}

func (e Event) String() string {
	return fmt.Sprintf("type: %s data: %s", e.Type, e.Data)
}
