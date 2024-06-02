package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

type configLine struct {
	hasAnything bool
	err         error
	operation   string
	arg1        string
	arg2        string
}

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
	if err := execCmd(true, "/usr/bin/go", "mod", "tidy"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	fmt.Println("Builiding init...")
	spinner.Start()
	if err := execCmd(true, "/usr/bin/go", "build", "-o", "mInit"); err != nil {
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
	if _, err := os.Stat(curDir + "/malino.cfg"); !os.IsNotExist(err) {
		file, err := os.ReadFile(curDir + "/malino.cfg")
		if err != nil {
			return err
		}
		lines := strings.Split(string(file), "\n")
		for _, line := range lines {
			confLine := parseConfigLine(line)
			if confLine.err != nil {
				return confLine.err
			}
			if err := handleLine(confLine); err != nil {
				return err
			}
		}
	}
	if err := execCmd(false, "/usr/bin/bash", "-c", "find . -print0 | cpio --null -ov --format=newc | gzip -9 > ../initramfs.cpio.gz"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	goToParentDir()
	if _, err := os.Stat("vmlinuz"); os.IsNotExist(err) {
		fmt.Println("Downloading kernel...")
		spinner.Start()
		getKernel()
		spinner.Stop()
	}

	if err := os.RemoveAll("initrd"); err != nil {
		return err
	}

	fmt.Println("Done.")

	return nil
}

func parseConfigLine(line string) configLine {
	// check if line is empty or is a comment, if it is, just say it didn't do anything
	if line == "" || strings.HasPrefix(line, "#") {
		return configLine{false, nil, "", "", ""}
	}
	// split by spaces, throw error if there is not 3 words since that's the syntax
	words := strings.Split(line, " ")
	if len(words) != 3 {
		return configLine{false, fmt.Errorf("line does not contain 3 words"), "", "", ""}
	}

	op := ""

	switch words[0] {
	case "include":
		op = "include"
	default:
		return configLine{false, fmt.Errorf("invalid operation"), "", "", ""}
	}

	return configLine{true, nil, op, words[1], words[2]}
}

func handleLine(line configLine) error {
	if !line.hasAnything {
		return nil
	}

	switch line.operation {
	case "include":
		fmt.Printf("including %v as %v in the malino system\n", line.arg1, line.arg2)
		curDir := "undefined"
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			curDir = dir
		}
		if strings.HasPrefix(line.arg1, "https://") {
			if err := downloadFile(line.arg1, "file_malinoAutoDownload.tmp"); err != nil {
				return err
			}
			if err := copy("file_malinoAutoDownload.tmp", curDir+line.arg2); err != nil {
				return err
			}
		}
		if err := copy(line.arg1, curDir+line.arg2); err != nil {
			return err
		}
	}

	return nil
}
