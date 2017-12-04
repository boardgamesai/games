package game

type Comms struct {
	EventLog *EventLog
}

func (c *Comms) NewEvents(order int) []Event {
	return c.EventLog.NewForPlayer(order)
}
