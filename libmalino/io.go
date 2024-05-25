package libmalino

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

var oldState *syscall.Termios

func init() {
	setNonCanonicalMode()
}

func UserLine() string {
	var buf [1]byte
	var cmdString strings.Builder
	for {
		n, err := syscall.Read(int(os.Stdin.Fd()), buf[:])
		if err != nil {
			fmt.Printf("Critical error while reading characters:\n%v", err)
			return ""
		}
		if n > 0 {
			char := buf[0]
			if char == '\n' {
				fmt.Println()
				return cmdString.String()
			} else if char == 127 { // ASCII code for backspace
				if cmdString.Len() > 0 {
					cmd := cmdString.String()
					if len(cmd) > 1 {
						cmdString.Reset()
						cmdString.WriteString(cmd[:len(cmd)-1])
						fmt.Print("\b \b")
					} else {
						cmdString.Reset()
						fmt.Print("\b \b")
					}
				}
			} else {
				fmt.Print(string(char))
				cmdString.WriteByte(char)
			}
		}
	}
}

func ClearScreen() {
	fmt.Print("\033[2J\033[H")
}

func setNonCanonicalMode() {
	fd := int(os.Stdin.Fd())
	var termios syscall.Termios
	_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCGETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if errno != 0 {
		fmt.Printf("Error getting terminal attributes: %v\n", errno)
		os.Exit(1)
	}
	oldState = &termios
	termios.Lflag &^= syscall.ICANON | syscall.ECHO
	termios.Cc[syscall.VMIN] = 1
	termios.Cc[syscall.VTIME] = 0
	_, _, errno = syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(&termios)), 0, 0, 0)
	if errno != 0 {
		fmt.Printf("Error setting terminal attributes: %v\n", errno)
		os.Exit(1)
	}
}

func resetTerminalMode() {
	if oldState != nil {
		fd := int(os.Stdin.Fd())
		_, _, errno := syscall.Syscall6(syscall.SYS_IOCTL, uintptr(fd), uintptr(syscall.TCSETS), uintptr(unsafe.Pointer(oldState)), 0, 0, 0)
		if errno != 0 {
			fmt.Printf("Error resetting terminal attributes: %v\n", errno)
		}
	}
}
