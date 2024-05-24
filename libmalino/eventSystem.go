package libmalino

type Event struct {
	ID        int
	Type      uint8
	TimeStamp string
	Data      string
	Caller    string
}

var Events []Event

func LogDebug(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 0, SystemUptimeAsString(), Data, Caller})
}

func LogInfo(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 1, SystemUptimeAsString(), Data, Caller})
}

func LogWarning(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 2, SystemUptimeAsString(), Data, Caller})
}

func LogError(Data string, Caller string) {
	Events = append(Events, Event{len(Events), 3, SystemUptimeAsString(), Data, Caller})
}

func LatestEvent() Event {
	return Events[len(Events)-1]
}

func AllEvents() []Event {
	return Events
}
