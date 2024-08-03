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

	name := "undefined"
	curDir := "undefined"
	if dir, err := os.Getwd(); err != nil {
		return err
	} else {
		curDir = dir
		name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
	}

	if _, err := os.Stat("malino.cfg"); os.IsNotExist(err) {
		return fmt.Errorf("current directory does not contain a valid malino project")
	}

	fmt.Println(" RM initramfs.cpio.gz")
	spinner.Start()
	if _, err := os.Stat("initramfs.cpio.gz"); !os.IsNotExist(err) {
		os.Remove("initramfs.cpio.gz")
	}

	var conf []configLine
	if _, err := os.Stat(curDir + "/malino.cfg"); !os.IsNotExist(err) {
		file, err := os.ReadFile(curDir + "/malino.cfg")
		if err != nil {
			return err
		}
		lines := strings.Split(string(file), "\n")
		for lineNum, line := range lines {
			confLine := parseConfigLine(line, lineNum+1)
			if confLine.err != nil {
				return confLine.err
			}
			if confLine.hasAnything {
				conf = append(conf, confLine)
			}
		}
	}

	/*for _, line := range conf {
		fmt.Printf("op: %v | args: %v\n", line.operation, strings.Join(line.args, " "))
	}*/

	lang, err := handleLangLine(conf[0])
	if err != nil {
		return err
	}
	switch lang {
	case "go":
		if err := buildGoProj(spinner, conf); err != nil {
			return err
		}
	case "c#":
		if err := buildCSProj(spinner, conf, name, curDir); err != nil {
			return err
		}
	default:
		return fmt.Errorf("malino.cfg: line 1: invalid language")
	}

	// TODO: compile other stuff
	// nah just let the user make a makefile
	// use maura as an example

	fmt.Println(" MK initramfs.cpio.gz")
	if err := createAndCD("initrd"); err != nil {
		spinner.Stop()
		return err
	}
	if err := os.Rename(curDir+"/mInit", curDir+"/initrd/init"); err != nil {
		os.RemoveAll("initrd")
		spinner.Stop()
		return err
	}
	for _, line := range conf {
		if line.operation == "include" {
			handleIncludeLine(line)
		}
	}
	if err := execCmd(false, "/usr/bin/bash", "-c", "find . -print0 | cpio --null -ov --format=newc | gzip -9 > ../initramfs.cpio.gz"); err != nil {
		os.RemoveAll("initrd")
		spinner.Stop()
		return err
	}

	goToParentDir()
	if _, err := os.Stat("/home/" + currentUser.Username + "/.malino/vmlinuz"); os.IsNotExist(err) {
		fmt.Println("Downloading Linux kernel & modules. This will take a long time (2-5 minutes), but you'll only have to do it once. (or when you run `malino update-kernel`)")
		spinner.Start()
		if err := getKernel(true); err != nil {
			os.RemoveAll("initrd")
			return err
		}
		spinner.Stop()
	}

	if err := os.RemoveAll("initrd"); err != nil {
		return err
	}

	// "malino export"
	exportProj(name)

	spinner.Stop()

	if _, err := os.Stat("vmlinuz"); err == nil {
		fmt.Println("btw, you can delete the vmlinuz file in this project, malino now stores it's kernel in ~/.malino/vmlinuz.")
	}

	return nil
}
