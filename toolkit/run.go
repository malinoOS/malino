package main

import (
	"fmt"
	"os"
)

func runProj() error {
	if _, err := os.Stat("bzImage"); os.IsNotExist(err) {
		return fmt.Errorf("current directory either doesn't contain a project or has not been built yet")
	}

	if _, err := os.Stat("initramfs.cpio.gz"); os.IsNotExist(err) {
		return fmt.Errorf("current directory either doesn't contain a project or has not been built yet")
	}

	if err := execCmd(true, "qemu-system-x86_64", "-kernel", "bzImage", "-initrd", "initramfs.cpio.gz"); err != nil {
		return err
	}

	fmt.Println("Done.")

	return nil
}
