package libmalino

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func ShutdownComputer() {
	fmt.Printf("syncing disks...\n")
	syscall.Sync()
	//fmt.Printf("unmounting disks...\n")
	//UnmountProcFS()
	fmt.Printf("shutting down...\n")
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

func RebootComputer() {
	fmt.Printf("syncing disks...\n")
	syscall.Sync()
	//fmt.Printf("unmounting disks...\n")
	//UnmountProcFS()
	fmt.Printf("shutting down...\n")
	syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

func SystemUptimeAsInt() int {
	dat, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0
	}
	stringRep := strings.Split(strings.Split(string(dat), " ")[1], ".")
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
	stringRep := strings.Split(string(dat), " ")[1]
	i, err := strconv.ParseFloat(stringRep, 64)
	if err != nil {
		return 0
	}
	return i
}

func SpawnProcess(path string, startDir string, environmentVariables []string, files []uintptr, wait bool, errorIfExit bool, args ...string) error {
	procAttr := &syscall.ProcAttr{
		Dir:   startDir,
		Env:   environmentVariables,
		Files: files,
		Sys:   nil,
	}
	var wstatus syscall.WaitStatus

	args = append(args, path)
	copy(args[1:], args)
	args[0] = path
	pid, err := syscall.ForkExec(path, args, procAttr)
	if err != nil {
		fmt.Printf("err: could not execute %v: %v\n", path, err.Error())
		return err
	} else {
		if wait {
			_, err = syscall.Wait4(pid, &wstatus, 0, nil)
			if err != nil {
				fmt.Printf("err: could not execute %v: %v\n", path, err.Error())
				return err
			}
		}
	}

	if wstatus.Exited() {
		if errorIfExit || wstatus.ExitStatus() != 0 {
			// Process exited
			// Create a new error
			LogError(fmt.Sprintf("%v exited with code %d", path, wstatus.ExitStatus()), "libmalino.SpawnProcess")
			return fmt.Errorf("%v exited with code %d", path, wstatus.ExitStatus())
		}
	}
	return nil
}

func SpawnProcessStdioFiles(path string, startDir string, environmentVariables []string, wait bool, errorIfExit bool, args ...string) error {
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
		return err
	} else {
		if wait {
			_, err = syscall.Wait4(pid, &wstatus, 0, nil)
			if err != nil {
				fmt.Printf("err: could not execute %v: %v\n", path, err.Error())
				return err
			}
		}
	}

	if wstatus.Exited() && errorIfExit && wstatus.ExitStatus() != 0 {
		// Process exited
		// Create a new error
		LogError(fmt.Sprintf("%v exited with code %d", path, wstatus.ExitStatus()), "libmalino.SpawnProcess")
		return fmt.Errorf("%v exited with code %d", path, wstatus.ExitStatus())
	}
	return nil
}

func MountProcFS() error {
	if err := os.Mkdir("/proc", 0777); err != nil {
		return err
	}
	if err := syscall.Mount("proc", "/proc", "proc", uintptr(0), ""); err != nil {
		return err
	}
	return nil
}

func UnmountProcFS() error {
	if err := syscall.Unmount("/proc", 0); err != nil {
		return err
	}
	return nil
}

func MountDevFS() error {
	if err := syscall.Mount("udev", "/dev", "devtmpfs", syscall.MS_NOSUID, ""); err != nil {
		return err
	}
	return nil
}

func UnmountDevFS() error {
	if err := syscall.Unmount("/dev", 0); err != nil {
		return err
	}
	return nil
}
