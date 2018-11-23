package game

import "encoding/json"

type Message struct {
	Type string
	Data json.RawMessage
}
