using System;

namespace libmalino;

/// <summary>
/// Functions for some general stuff relating to input/output from/to the the user.
/// </summary>
public class malinoIO {
    /// <summary>
    /// Clears the screen.
    /// </summary>
    public static void ClearScreen() {
        Console.Write("\x1b[2J\x1b[H");
    }
}