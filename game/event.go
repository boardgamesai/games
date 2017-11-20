package game

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	Type string
	Data json.RawMessage
	Show map[int]bool `json:"-"` // Player orders who should be shown this event
	Seen map[int]bool `json:"-"` // Player orders who have seen this event
}

func (e Event) String() string {
	return fmt.Sprintf("type: %s data: %s", e.Type, e.Data)
}
