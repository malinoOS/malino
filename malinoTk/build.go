package main

import (
	"fmt"
	"os"
	"os/exec"
)

func buildProj() error {
	println("checking if go.mod exists...")
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory does not contain a valid malino project")
	}

	println("getting dependencies...")
	cmd := exec.Command("/usr/bin/go", "mod", "tidy")
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))

	println("building project...")
	cmd = exec.Command("/usr/bin/go", "build", "-o", "malinoOS")
	stdout, err = cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))

	println("creating makefile...")
	err = os.WriteFile("Makefile", []byte(
		"all:\n"+
			"	sudo modprobe nbd max_part=8\n"+
			"	sudo qemu-nbd -c /dev/$(DRIVE) golinux-main/linux.qcow2\n"+
			"	mkdir -p disk\n"+
			"	sudo mount -t ext4 /dev/$(DRIVE)p1 disk\n"+
			"	sudo mv malinoOS disk/sbin/malino\n"+
			"	sudo umount disk\n"+
			"	rm -rf disk\n"+
			"	sudo qemu-nbd -d /dev/$(DRIVE)\n"), 0777)
	if err != nil {
		return err
	}

	println("running make...")
	cmd = exec.Command("/usr/bin/make")
	stdout, err = cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))

	return nil
}