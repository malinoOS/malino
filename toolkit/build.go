package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func buildProj() error {
	// Initialize the spinner (loading thing).
	spinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory does not contain a valid malino project")
	}

	fmt.Println("Removing binaries...")
	spinner.Start()
	if _, err := os.Stat("initramfs.cpio.gz"); !os.IsNotExist(err) {
		os.Remove("initramfs.cpio.gz")
	}
	spinner.Stop()

	name := "undefined"
	curDir := "undefined"
	if dir, err := os.Getwd(); err != nil {

		return err
	} else {
		name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
		curDir = dir
	}
	fmt.Println("Found name: " + name)

	fmt.Println("Getting dependencies...")
	spinner.Start()
	if err := execCmd(false, "/usr/bin/go", "mod", "tidy"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	fmt.Println("Builiding init...")
	spinner.Start()
	if err := execCmd(false, "/usr/bin/go", "build", "-o", "mInit"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	// TODO: compile other stuff

	fmt.Println("Creating initramfs...")
	spinner.Start()
	if err := createAndCD("initrd"); err != nil {
		spinner.Stop()
		return err
	}
	if err := os.Rename(curDir+"/mInit", curDir+"/initrd/init"); err != nil {
		spinner.Stop()
		return err
	}
	if err := execCmd(false, "/usr/bin/bash", "-c", "find . -print0 | cpio --null -ov --format=newc | gzip -9 > ../initramfs.cpio.gz"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	goToParentDir()
	if _, err := os.Stat("bzImage"); os.IsNotExist(err) {
		fmt.Println("Downloading kernel...")
		spinner.Start()
		if err := downloadFile("https://winksplorer.net/bzImage", "bzImage"); err != nil {
			spinner.Stop()
			return err
		}
		spinner.Stop()
	}

	if err := os.RemoveAll("initrd"); err != nil {
		return err
	}

	fmt.Println("Done.")

	return nil
}
