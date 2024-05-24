package main

import (
	"fmt"
	"os"
	"os/exec"
)

func runProj() error {
	println("Checking if go.mod exists...")
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory does not contain a valid malino project")
	}

	println("Starting QEMU...")
	cmd := exec.Command("qemu-system-x86_64", "-drive", "file=golinux-main/linux.qcow2,format=qcow2", "-m", "4G", "-enable-kvm", "-smp", "4")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	fmt.Println(string(stdout))

	return nil
}
