package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func exportProj() error {
	// Initialize the spinner (loading thing).
	spinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	name := "undefined"
	if dir, err := os.Getwd(); err != nil {
		return err
	} else {
		name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
	}
	fmt.Println("Found name: " + name)

	if _, err := os.Stat("bzImage"); os.IsNotExist(err) {
		return fmt.Errorf("current directory either doesn't contain a project or has not been built yet")
	}

	if _, err := os.Stat("initramfs.cpio.gz"); os.IsNotExist(err) {
		return fmt.Errorf("current directory either doesn't contain a project or has not been built yet")
	}

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory doesn't contain a project")
	}

	fmt.Println("Creating folders...")
	spinner.Start()
	if err := createAndCD("iso"); err != nil {
		spinner.Stop()
		return err
	}
	if err := createAndCD("boot"); err != nil {
		spinner.Stop()
		return err
	}
	if err := createAndCD("grub"); err != nil {
		spinner.Stop()
		return err
	}
	goToParentDir()
	goToParentDir()
	goToParentDir()
	spinner.Stop()

	fmt.Println("Moving files...")
	spinner.Start()
	if err := os.Rename("bzImage", "iso/boot/bzImage"); err != nil {
		spinner.Stop()
		return err
	}
	if err := os.Rename("initramfs.cpio.gz", "iso/boot/initramfs.cpio.gz"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	fmt.Println("Running grub...")
	spinner.Start()
	err := os.WriteFile("iso/boot/grub/grub.cfg", []byte(
		"set default=0\n"+
			"set timeout=0\n\n"+
			"menuentry '"+name+"' --class os {\n"+
			"	insmod gzio\n"+
			"	insmod part_msdos\n"+
			"	linux /boot/bzImage\n"+
			"	initrd /boot/initramfs.cpio.gz\n"+
			"}"), 0777)
	if err != nil {
		spinner.Stop()
		return err
	}

	if err := execCmd(false, "grub-mkrescue", "-o", name+".iso", "iso/"); err != nil {
		spinner.Stop()
		return err
	}

	if err := os.RemoveAll("iso"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	fmt.Println("Done.")

	return nil
}
