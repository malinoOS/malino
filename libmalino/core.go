package libmalino

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func ShutdownComputer() error {
	fmt.Printf("syncing disks...\n")
	syscall.Sync()
	//fmt.Printf("unmounting disks...\n")
	//UnmountProcFS()
	fmt.Printf("shutting down...\n")
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func RebootComputer() error {
	fmt.Printf("syncing disks...\n")
	syscall.Sync()
	//fmt.Printf("unmounting disks...\n")
	//UnmountProcFS()
	fmt.Printf("shutting down...\n")
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func SystemUptimeAsInt() int {
	dat, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	stringRep := strings.Split(strings.Split(string(dat), " ")[0], ".")
	i, err := strconv.Atoi(stringRep[0])
	if err != nil {
		return 0
	}
	return i
}

func SystemUptimeAsFloat() float64 {
	dat, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	stringRep := strings.Split(string(dat), " ")[0]
	i, err := strconv.ParseFloat(stringRep, 64)
	if err != nil {
		return 0
	}
	return i
}

func SpawnProcess(path string, startDir string, environmentVariables []string, wait bool, args ...string) (int, error) {
	procAttr := &syscall.ProcAttr{
		Dir:   startDir,
		Env:   environmentVariables,
		Files: []uintptr{os.Stdout.Fd(), os.Stdin.Fd(), os.Stderr.Fd()},
		Sys:   nil,
	}
	var wstatus syscall.WaitStatus

	args = append(args, path)
	copy(args[1:], args)
	args[0] = path
	pid, err := syscall.ForkExec(path, args, procAttr)
	if err != nil {
		fmt.Printf("err: could not execute %v: %v\n", path, err.Error())
		return -1, err
	}

	if wait {
		_, err = syscall.Wait4(pid, &wstatus, 0, nil)
		if err != nil {
			fmt.Printf("err: could not execute %v: %v\n", path, err.Error())
			return -1, err
		}

		if wstatus.Exited() {
			return wstatus.ExitStatus(), nil
		}
	}

	return 0, nil
}

func MountProcFS() error {
	if err := os.Mkdir("/proc", 0777); err != nil {
		return err
	}
	return syscall.Mount("proc", "/proc", "proc", uintptr(0), "")
}

func UnmountProcFS() error {
	return syscall.Unmount("/proc", 0)
}

func MountDevFS() error {
	return syscall.Mount("udev", "/dev", "devtmpfs", syscall.MS_NOSUID, "")
}

func UnmountDevFS() error {
	return syscall.Unmount("/dev", 0)
}
