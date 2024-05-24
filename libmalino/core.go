package libmalino

import (
	"fmt"
	"strconv"
	"strings"
)

func ShutdownComputer() {
	fmt.Printf("syncing disks...\n")
	syscall.Sync()
	fmt.Printf("shutting down...\n")
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func SystemUptimeAsInt() int {
	dat, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	stringRep := strings.Split(strings.Split(string(dat), " ")[1], ".")
	i, err := strconv.Atoi(stringRep)
	if err != nil {
		return 0
	}
	return i
}

func SystemUptimeAsString() string {
	dat, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	return strings.Split(strings.Split(string(dat), " ")[1], ".")
}
