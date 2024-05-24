package libmalino

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"
)

func UserLine() string {
	var buf [1]byte
	var cmdString strings.Builder
	for {
		n, err := syscall.Read(int(os.Stdin.Fd()), buf[:])
		if err != nil {
			fmt.Printf("Critcal error while reading characters:\n%v", err)
			for true {
			}
		}
		if n > 0 {
			char := buf[0]
			if char == '\n' {
				fmt.Println()
				return cmdString.String()
			} else if char == 127 { // ASCII code for backspace
				if cmdString.Len() > 0 {
					// Convert the builder to a string, remove the last character, and create a new builder
					cmd := cmdString.String()
					if len(cmd) > 1 {
						cmdString.Reset()
						cmdString.WriteString(cmd[:len(cmd)-1])
						// Move the cursor back and clear the character
						fmt.Print("\b \b")
					} else {
						// If only one character is present, simply reset the builder
						cmdString.Reset()
						fmt.Print("\b \b") // Clear the character
					}
				}
			} else {
				fmt.Print(string(char))
				cmdString.WriteByte(char)
			}
		} else {
			time.Sleep(time.Millisecond * 10)
		}
	}
}
