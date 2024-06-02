package main

import (
	"fmt"
	"os"
)

func runProj(args []string) error {
	if _, err := os.Stat("vmlinuz"); os.IsNotExist(err) {
		if err := buildProj(); err != nil {
			return err
		}
	}

	if _, err := os.Stat("initramfs.cpio.gz"); os.IsNotExist(err) {
		if err := buildProj(); err != nil {
			return err
		}
	}

	hasSerialArg := len(args) == 2
	if hasSerialArg && args[1] == "-serial" {
		if err := execCmdDirectStdio("qemu-system-x86_64", "-kernel", "vmlinuz", "-initrd", "initramfs.cpio.gz", "-append", "console=ttyS0", "-nographic", "-m", "512M", "-enable-kvm"); err != nil {
			return err
		}
	} else {
		if err := execCmd(true, "qemu-system-x86_64", "-kernel", "vmlinuz", "-initrd", "initramfs.cpio.gz", "-m", "512M", "-enable-kvm"); err != nil {
			return err
		}
	}

	fmt.Println("Done.")

	return nil
}
