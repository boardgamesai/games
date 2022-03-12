package game

import "encoding/json"

type Message struct {
	Type string
	Data json.RawMessage
}

type MessageResponse struct {
	Err  *DQError `json:",omitempty"`
	Data json.RawMessage
}
