package libmalino

// An event.
//
// ID = Identification number for the event. Stored in chronological order automatically.
//
// Type = The type of the event. 0 = debug, 1 = info, 2 = warning, 3 = error.
//
// TimeStamp = When the event happened. This is handled automatically. See SystemUptimeAsFloat().
//
// Data = The actual event message.
//
// Caller = Where the event originated from.
type Event struct {
	ID        int
	Type      uint8
	TimeStamp float64
	Data      string
	Caller    string
}

// The global event slice. Stores all known events that your OS made.
var Events []Event

// Logs an event to the global event slice as type debug.
func LogDebug(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 0, SystemUptimeAsFloat(), Data, Caller})
}

// Logs an event to the global event slice as type info.
func LogInfo(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 1, SystemUptimeAsFloat(), Data, Caller})
}

// Logs an event to the global event slice as type warning.
func LogWarning(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 2, SystemUptimeAsFloat(), Data, Caller})
}

// Logs an event to the global event slice as type error.
func LogError(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 3, SystemUptimeAsFloat(), Data, Caller})
}

// Returns: the latest event in the global event slice.
func LatestEvent() Event {
	return Events[len(Events)-1]
}

// Returns: the global event slice.
func AllEvents() []Event {
	return Events
}
