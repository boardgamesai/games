package game

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var AllPlayers = []int{-1}

type EventLog []Event

func (el *EventLog) Add(event interface{}, orders []int) error {
	eventType := reflect.TypeOf(event).Name()
	if eventType[0:5] != "Event" {
		return fmt.Errorf("Invalid type %s passed to AddEvent", eventType)
	}
	eventType = strings.ToLower(eventType[5:])

	eventJSON, err := json.Marshal(&event)
	if err != nil {
		return err
	}

	show := map[int]bool{}
	for _, order := range orders {
		show[order] = true
	}

	e := Event{
		Type: eventType,
		Data: eventJSON,
		Show: show,
		Seen: map[int]bool{},
	}
	*el = append(*el, e)

	return nil
}

func (el *EventLog) NewForPlayer(order int) []Event {
	events := []Event{}

	// We start at the most recent event and move backwards until we find one they've already seen.
	for i := len(*el) - 1; i >= 0; i-- {
		e := &((*el)[i])

		if e.Seen[order] {
			// Found an already-seen one, we're done.
			break
		}

		// OK, this is new, mark it seen.
		e.Seen[order] = true

		// But should we include it?
		if !(e.Show[-1] || e.Show[order]) {
			continue
		}

		// If we're here, they get this event, but we're going backwards so append args get reversed.
		events = append([]Event{*e}, events...)
	}

	return events
}

func (el *EventLog) Clear() {
	*el = []Event{}
}
