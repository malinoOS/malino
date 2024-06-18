using System.Collections.Generic;

namespace libmalino;

public struct Event {
	public int ID;
    public EventType Type;
    public float TimeStamp;
    public string Data;
    public string Caller;
}

public enum EventType {
    Debug,
    Info,
    Warning,
    Error
}

public class EventLogger {
    public static List<Event> Events = new();

    public static void LogDebug(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Debug,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    public static void LogInfo(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Info,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    public static void LogWarning(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Warning,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    public static void LogError(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Error,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    public static Event LatestEvent() {
        return Events[Events.Count-1];
    }

    public static Event[] AllEvents() {
        return Events.ToArray();
    }
}