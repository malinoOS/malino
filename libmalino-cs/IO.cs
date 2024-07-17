using System;

namespace libmalino;

/// <summary>
/// Functions for some general stuff relating to input/output from/to the the user.
/// </summary>
public class malinoIO {
    /// <summary>
    /// Clears the screen. Because Console.Clear() doesn't work for some reason.
    /// </summary>
    public static void ClearScreen() {
        Console.Write("\x1b[2J\x1b[H");
    }

    /// <summary>
    /// Reads a line from the user, always use this function instead of C#'s built-in way, because backspace
    /// </summary>
    public static string UserLine()
    {
        string input = "";
        while (true)
        {
            ConsoleKeyInfo key = Console.ReadKey(true); // read a byte

            if (key.Key == ConsoleKey.Enter) // newline
            {
                Console.WriteLine();
                break;
            }
            else if (key.Key == ConsoleKey.Backspace) // backspace
            {
                if (input.Length > 0)
                {
                    input = input.Remove(input.Length - 1);
                    Console.Write("\b \b"); // Clears the character visually on the screen
                }
            }
            else // anything else
            {
                input += key.KeyChar;
                Console.Write(key.KeyChar);
            }
        }
        return input;
    }
}