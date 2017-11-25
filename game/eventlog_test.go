package game

import (
	"fmt"
	"strings"
	"testing"
)

type EventTest1 struct {
	Val int
}

type EventTest2 struct {
	Val int
}

type EventTest3 struct {
	Val int
}

func getEventLog() (*EventLog, error) {
	l := &EventLog{}

	var err error
	if err = l.Add(EventTest1{Val: 1}, AllPlayers); err != nil {
		return l, err
	}
	if err = l.Add(EventTest2{Val: 2}, []int{2}); err != nil {
		return l, err
	}
	if err = l.Add(EventTest3{Val: 3}, []int{1}); err != nil {
		return l, err
	}

	return l, nil
}

func TestAddEventType(t *testing.T) {
	l, err := getEventLog()
	if err != nil {
		t.Fatalf("Adding events returned error: %s", err)
	}

	for i, e := range *l {
		typeName := fmt.Sprintf("test%d", i+1)
		if e.Type != typeName {
			t.Errorf("Didn't get type %s for event, got: %s", typeName, e.Type)
		}
	}
}

func TestAddNonEvent(t *testing.T) {
	l := &EventLog{}
	r := Place{} // Arbitrary, just can't be an Event*
	err := l.Add(r, AllPlayers)
	if err == nil {
		t.Errorf("Expected error adding non-event, didn't get one")
	} else if !strings.Contains(err.Error(), "Invalid type") {
		t.Errorf("Didn't find expected error, instead got error: %s", err)
	}
}

func TestNewForPlayer(t *testing.T) {
	l, err := getEventLog()
	if err != nil {
		t.Fatalf("Adding events returned error: %s", err)
	}

	tests1 := []struct {
		order    int
		expected []int
	}{
		{1, []int{1, 3}},
		{2, []int{1, 2}},
		{3, []int{1}},
	}

	for _, test := range tests1 {
		events, ok := checkLog(l, test.order, test.expected)
		if !ok {
			t.Errorf("Found unexpected events for %d: %+v", test.order, events)
		}
	}

	l.Add(EventTest3{Val: 4}, []int{2})
	l.Add(EventTest2{Val: 5}, AllPlayers)
	l.Add(EventTest1{Val: 6}, []int{2, 3})
	l.Add(EventTest3{Val: 7}, AllPlayers)
	l.Add(EventTest2{Val: 8}, []int{2})
	l.Add(EventTest1{Val: 9}, []int{1})

	tests2 := []struct {
		order    int
		expected []int
	}{
		{1, []int{2, 3, 1}},
		{2, []int{3, 2, 1, 3, 2}},
		{3, []int{2, 1, 3}},
		{4, []int{1, 2, 3}}, // The 1 is from getEventLog()
	}

	for _, test := range tests2 {
		events, ok := checkLog(l, test.order, test.expected)
		if !ok {
			t.Errorf("Found unexpected events for %d: %+v", test.order, events)
		}
	}
}

func checkLog(l *EventLog, order int, expected []int) ([]Event, bool) {
	events := l.NewForPlayer(order)
	if len(events) != len(expected) {
		return events, false
	}

	for i := 0; i < len(events); i++ {
		typeName := fmt.Sprintf("test%d", expected[i])
		if events[i].Type != typeName {
			return events, false
		}
	}

	return events, true
}
