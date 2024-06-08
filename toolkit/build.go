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
	args        []string
}

func buildProj() error {
	// Initialize the spinner (loading thing).
	spinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory does not contain a valid malino project")
	}

	fmt.Println(" RM initramfs.cpio.gz")
	spinner.Start()
	if _, err := os.Stat("initramfs.cpio.gz"); !os.IsNotExist(err) {
		os.Remove("initramfs.cpio.gz")
	}
	spinner.Stop()

	//name := "undefined"
	curDir := "undefined"
	if dir, err := os.Getwd(); err != nil {

		return err
	} else {
		//name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
		curDir = dir
	}

	var conf []configLine
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
			if confLine.hasAnything {
				conf = append(conf, confLine)
			}
		}
	}

	/*for _, line := range conf {
		fmt.Printf("op: %v | args: %v\n", line.operation, strings.Join(line.args, " "))
	}*/

	fmt.Println(" DL dependencies")
	spinner.Start()
	if err := execCmd(true, "/usr/bin/go", "mod", "tidy"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	fmt.Println(" GO init")
	spinner.Start()
	buildFlagsExist := false
	for _, line := range conf {
		if line.operation == "buildflags" {
			buildFlagsExist = true
			if err := execCmd(true, append([]string{"/usr/bin/go", "build", "-o", "mInit"}, line.args...)...); err != nil {
				spinner.Stop()
				return err
			}
		} else if line.operation == "verfmt" {
			buildFlagsExist = true
			ver, err := handleVerfmtLine(line)
			if err != nil {
				return err
			}
			if err := execCmd(true, "/usr/bin/go", "build", "-o", "mInit", "-ldflags", "-X main.Version="+ver); err != nil {
				spinner.Stop()
				return err
			}
		}
	}
	if !buildFlagsExist {
		if err := execCmd(true, "/usr/bin/go", "build", "-o", "mInit"); err != nil {
			spinner.Stop()
			return err
		}
	}
	spinner.Stop()

	// TODO: compile other stuff
	// nah just let the user make a makefile
	// use maura as an example

	fmt.Println(" MK initramfs.cpio.gz")
	spinner.Start()
	if err := createAndCD("initrd"); err != nil {
		spinner.Stop()
		return err
	}
	if err := os.Rename(curDir+"/mInit", curDir+"/initrd/init"); err != nil {
		spinner.Stop()
		return err
	}
	for _, line := range conf {
		if line.operation == "include" {
			handleIncludeLine(line)
		}
	}
	if err := execCmd(false, "/usr/bin/bash", "-c", "find . -print0 | cpio --null -ov --format=newc | gzip -9 > ../initramfs.cpio.gz"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	goToParentDir()
	if _, err := os.Stat("vmlinuz"); os.IsNotExist(err) {
		fmt.Println(" DL vmlinuz")
		spinner.Start()
		getKernel()
		spinner.Stop()
	}

	if err := os.RemoveAll("initrd"); err != nil {
		return err
	}

	return nil
}

func parseConfigLine(line string) configLine {
	// check if line is empty or is a comment, if it is, just say it didn't do anything
	if line == "" || strings.HasPrefix(line, "#") {
		return configLine{false, nil, "", nil}
	}
	// split by spaces, throw error if there is not 3 words since that's the syntax
	words := strings.Split(line, " ")

	op := ""

	switch words[0] {
	case "include":
		if len(words) != 3 {
			return configLine{false, fmt.Errorf("line does not contain 3 words, which is required for include operation"), "", nil}
		}
		op = "include"
	case "buildflags":
		return configLine{true, nil, "buildflags", combineQuotedStrings(words[1:])}
	case "verfmt":
		if len(words) != 2 {
			return configLine{false, fmt.Errorf("line does not contain 2 words, which is required for verfmt operation"), "", nil}
		}
		op = "verfmt"
	//case "lang":
	//	op = "lang"
	default:
		return configLine{false, fmt.Errorf("invalid operation"), "", nil}
	}

	return configLine{true, nil, op, words[1:]}
}

func handleIncludeLine(line configLine) error {
	if !line.hasAnything {
		return fmt.Errorf("the entire configuration parser is broken. good luck")
	}

	if line.operation == "include" {
		fmt.Printf("INC %v AS %v\n", line.args[0], line.args[1])
		curDir := "undefined"
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			curDir = dir
		}
		line.args[0] = strings.Replace(line.args[0], ".", curDir, 1)
		if strings.HasPrefix(line.args[0], "https://") {
			if err := downloadFile(line.args[0], "file_malinoAutoDownload.tmp"); err != nil {
				return err
			}
			if err := copyFile("file_malinoAutoDownload.tmp", curDir+line.args[1]); err != nil {
				return err
			}
			if err := os.Remove("file_malinoAutoDownload.tmp"); err != nil {
				return err
			}
			return nil
		}
		if strings.HasPrefix(line.args[0], "dir///") {
			if err := copyDirectory(line.args[0][6:], curDir+line.args[1]); err != nil {
				return err
			}
			return nil
		}
		if err := copyFile(line.args[0], curDir+line.args[1]); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("include handler called for non-include operation")
	}

	return nil
}

func handleVerfmtLine(line configLine) (string, error) {
	if !line.hasAnything {
		return "", fmt.Errorf("the entire configuration parser is broken. good luck")
	}

	if line.operation == "verfmt" {
		fmt.Printf("VERSION FORMAT: %v\n", line.args[0])
		switch line.args[0] {
		case "yymmdd":
			return time.Now().Format("060102"), nil
		case "ddmmyy":
			return time.Now().Format("020106"), nil
		case "mmddyy":
			return time.Now().Format("010206"), nil
		}
	} else {
		return "", fmt.Errorf("verfmt handler called for non-verfmt operation")
	}
	return "", fmt.Errorf("verfmt did nothing")
}
