using System.Collections.Generic;

namespace libmalino;

/// <summary>
/// An event.
/// </summary>
public struct Event {
    /// <summary>
    /// Identification number for the event. Stored in chronological order automatically.
    /// </summary>
	public int ID;
    /// <summary>
    /// The type of the event. Stored as <see cref="EventType"> libmalino.EventType</see>.
    /// </summary>
    public EventType Type;
    /// <summary>
    /// When the event happened. This is handled automatically. See <see cref="malino.SystemUptimeAsFloat">SystemUptimeAsFloat()</see>.
    /// </summary>
    public float TimeStamp;
    /// <summary>
    /// The actual event message.
    /// </summary>
    public string Data;
    /// <summary>
    /// Where the event originated from.
    /// </summary>
    public string Caller;
}

/// <summary>
/// The type of event.
/// </summary>
#pragma warning disable CS1591
public enum EventType {
    Debug,
    Info,
    Warning,
    Error
}
#pragma warning restore CS1591

/// <summary>
/// An event logger.
/// </summary>
public class EventLogger {
    /// <summary>
    /// The global event slice. Stores all known events that your OS made.
    /// </summary>
    public static List<Event> Events = new();

    /// <summary>
    /// Logs an event to the global event slice as type debug.
    /// </summary>
    public static void LogDebug(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Debug,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    /// <summary>
    /// Logs an event to the global event slice as type info.
    /// </summary>
    public static void LogInfo(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Info,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    /// <summary>
    /// Logs an event to the global event slice as type warning.
    /// </summary>
    public static void LogWarning(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Warning,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    /// <summary>
    /// Logs an event to the global event slice as type error.
    /// </summary>
    public static void LogError(string data, string caller) {
        Events.Add(new Event {
            ID = Events.Count,
            Type = EventType.Error,
            TimeStamp = malino.SystemUptimeAsFloat(),
            Data = data,
            Caller = caller
        });
    }

    /// <summary>
    /// Returns the latest event in the global event slice.
    /// </summary>
    public static Event LatestEvent() {
        return Events[Events.Count-1];
    }

    /// <summary>
    /// Returns the global event slice.
    /// </summary>
    public static Event[] AllEvents() {
        return Events.ToArray();
    }
}