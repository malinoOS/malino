package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/briandowns/spinner"
)

func buildProj() error {
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	s.Start()
	println("Checking if project exists...")
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory does not contain a valid malino project")
	}
	s.Stop()

	println("Getting dependencies...")
	s.Start()
	cmd := exec.Command("/usr/bin/go", "mod", "tidy")
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	s.Stop()

	println("Building project...")
	s.Start()
	cmd = exec.Command("/usr/bin/go", "build", "-o", "malinoOS")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	fmt.Println(string(stdout))
	s.Stop()

	println("Creating makefile...")
	s.Start()
	err = os.WriteFile("Makefile", []byte(
		"all:\n"+
			"ifeq ($(MALINO_KEY), malinoGoodCosmosBad)\n"+
			"	sudo modprobe nbd max_part=8\n"+
			"	sudo qemu-nbd -c /dev/nbd0 golinux-main/linux.qcow2\n"+
			"	mkdir -p disk\n"+
			"	sudo mount -t ext4 /dev/nbd0p1 disk\n"+
			"	sudo mv malinoOS disk/sbin/malino\n"+
			"	sudo umount disk\n"+
			"	rm -rf disk\n"+
			"	sudo qemu-nbd -d /dev/nbd0\n"+
			"else\n"+
			"	$(error This environment isn't the malino builder)\n"+
			"endif\n"), 0777)
	if err != nil {
		return err
	}
	s.Stop()

	println("Running make...")
	s.Start()
	cmd = exec.Command("/usr/bin/make", "MALINO_KEY=malinoGoodCosmosBad")
	stdout, err = cmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	s.Stop()

	println("Deleting makefile...")
	s.Start()
	err = os.Remove("Makefile")
	if err != nil {
		return err
	}
	s.Stop()

	return nil
}
