package libmalino

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
	"unsafe"
)

// Syncs disks then shuts down.
//
// Returns: Error or nil
func ShutdownComputer() error {
	// sync disks & shutdown
	fmt.Printf("syncing disks...\n")
	syscall.Sync()
	fmt.Printf("shutting down...\n")
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
}

// Reboots the computer.
//
// Returns: Error or nil
func RebootComputer() error {
	// sync disks & reboot
	fmt.Printf("syncing disks...\n")
	syscall.Sync()
	fmt.Printf("shutting down...\n")
	return syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
}

// Returns: The amount of seconds since Linux has booted. If it failed, it will return -1.
func SystemUptimeAsInt() int {
	// read uptime file
	dat, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return -1
	}
	// get the first number, without decimals
	stringRep := strings.Split(strings.Split(string(dat), " ")[0], ".")
	// parse as int
	i, err := strconv.Atoi(stringRep[0])
	if err != nil {
		return -1
	}
	return i
}

// Returns: The amount of seconds to the hundredths since Linux has booted. If it failed, it will return -1.
func SystemUptimeAsFloat() float64 {
	// read uptime file
	dat, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return -1
	}
	// get the first number
	stringRep := strings.Split(string(dat), " ")[0]
	// parse as float64
	i, err := strconv.ParseFloat(stringRep, 64)
	if err != nil {
		return -1
	}
	return i
}

// Spawns a process on the system.
//
// Returns: if wait=true, (return code, error), if wait=false, (0, error). Return code will be -1 if an error happened while spawning the process.
func SpawnProcess(path string, startDir string, environmentVariables []string, wait bool, args ...string) (int, error) {
	// get attributes
	procAttr := &syscall.ProcAttr{
		Dir:   startDir,
		Env:   environmentVariables,
		Files: []uintptr{os.Stdout.Fd(), os.Stdin.Fd(), os.Stderr.Fd()},
		Sys:   nil,
	}
	var wstatus syscall.WaitStatus

	// make args[0] be the path, as that is what it expects, and to make sense to the user
	args = append(args, path)
	copy(args[1:], args)
	args[0] = path

	// fork & execute: spawn a NEW process
	pid, err := syscall.ForkExec(path, args, procAttr)
	if err != nil {
		fmt.Printf("err: could not execute %v: %v\n", path, err.Error())
		return -1, err
	}

	// if wait, then wait & return the exit status
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

// Mounts /proc. /proc is recommended if you want your OS to like do stuff.
func MountProcFS() error {
	if err := os.Mkdir("/proc", 0777); err != nil {
		return err
	}
	return syscall.Mount("proc", "/proc", "proc", uintptr(0), "")
}

// Unmounts /proc.
func UnmountProcFS() error {
	return syscall.Unmount("/proc", 0)
}

// Mounts /dev. /dev is also recommended if you want your OS to like do stuff.
func MountDevFS() error {
	return syscall.Mount("udev", "/dev", "devtmpfs", syscall.MS_NOSUID, "")
}

// Unmounts /dev.
func UnmountDevFS() error {
	return syscall.Unmount("/dev", 0)
}

// Loads a Linux kernel module (.ko format).
func LoadKernelModule(modulePath string, params string) error {
	// open the module file
	fd, err := os.Open(modulePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	// get the size
	fileInfo, err := fd.Stat()
	if err != nil {
		return err
	}
	imageSize := fileInfo.Size()

	// alloc memory for the module
	image := make([]byte, imageSize)

	// read the module
	_, err = fd.Read(image)
	if err != nil {
		return err
	}

	byteptr, err := syscall.BytePtrFromString(params)
	if err != nil {
		return err
	}

	// syscall
	_, _, errno := syscall.Syscall(syscall.SYS_INIT_MODULE, uintptr(unsafe.Pointer(&image[0])), uintptr(imageSize), uintptr(unsafe.Pointer(byteptr)))
	if errno != 0 {
		return errno
	}

	return nil
}
