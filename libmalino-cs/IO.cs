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

    /// <summary>
    /// Reads a line from the user, always use this function instead of C#'s built-in way, because backspace
    /// </summary>
    public static string UserLine()
    {
        string input = "";
        while (true)
        {
            ConsoleKeyInfo key = Console.ReadKey(true);

            if (key.Key == ConsoleKey.Enter)
            {
                Console.WriteLine();
                break;
            }
            else if (key.Key == ConsoleKey.Backspace)
            {
                if (input.Length > 0)
                {
                    input = input.Remove(input.Length - 1);
                    Console.Write("\b \b"); // Clears the character visually on the screen
                }
            }
            else
            {
                input += key.KeyChar;
                Console.Write(key.KeyChar);
            }
        }
        return input;
    }
}