package game

type Comms struct {
	EventLog *EventLog
}

func (c *Comms) NewEvents(playerID PlayerID) []Event {
	return c.EventLog.NewForPlayer(playerID)
}
