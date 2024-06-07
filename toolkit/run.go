package main

import (
	"fmt"
	"os"
	"strings"
)

func runProj() error {
	name := "undefined"
	if dir, err := os.Getwd(); err != nil {
		return err
	} else {
		name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
	}

	if _, err := os.Stat(name + ".iso"); os.IsNotExist(err) {
		if err := exportProj([]string{args[0], "-efi"}); err != nil {
			return err
		}
	}

	if err := execCmd(true, "qemu-system-x86_64", "-m", "1G", "-enable-kvm", "-cdrom", name+".iso"); err != nil {
		return err
	}

	fmt.Println("Done.")

	return nil
}
