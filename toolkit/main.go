package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var args []string
var Version string = "undefined"

func main() {
	args = os.Args[1:]

	if len(args) > 0 { // check if there is any args after "malino"
		switch args[0] { // if there is, then check what it is
		case "help":
			printHelp()
			os.Exit(0)
		case "new":
			if err := newProj(args); err != nil {
				fmt.Printf("Error while creating project: %v", err.Error())
			}
		case "build":
			if err := buildProj(); err != nil {
				fmt.Printf("Error while building project: %v", err.Error())
			}
		case "run":
			if err := runProj(args); err != nil {
				fmt.Printf("Error while running project: %v", err.Error())
			}
		case "export":
			if err := exportProj(args); err != nil {
				fmt.Printf("Error while exporting project: %v", err.Error())
			}
		case "download-kernel":
			if err := getKernel(); err != nil {
				fmt.Printf("Error while downloading kernel: %v", err.Error())
			}
		default:
			fmt.Println("malino: Invalid operation")
			printHelp()
			os.Exit(1)
		}
	} else {
		fmt.Println("malino: No operation")
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Print(
		"malino toolkit v" + Version + "\n\n" +
			"malino help        	 Shows this help menu\n" +
			"malino new [name]  	 New project, creates folder and go module with name [name]\n" +
			"malino new         	 New project, does not make folder, and uses the name of the folder it's executed in\n" +
			"malino build       	 Builds a cpio of your OS\n" +
			"malino run            	 Runs your OS cpio with a precompiled linux\n" +
			"malino run -serial      Runs your OS cpio with a precompiled linux, but no qemu window shows, and interacts in stdio\n" +
			"malino export           Exports your OS into a .ISO file which can be shared or ran on real hardware BIOS machines\n" +
			"malino export -efi		 Exports your OS into an EFI .ISO file which can be shared or ran on real hardware UEFI machines\n" +
			"malino download-kernel  Downloads the latest Ubuntu Linux kernel.\n")
}

func createAndCD(dir string) error {
	// Create a directory and CD into it.
	err := os.Mkdir(dir, 0777)
	if err != nil {
		return err
	}

	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	return nil
}

func goToParentDir() {
	currentDir, err := os.Getwd()
	if err != nil {
		// give up, if you can't do a cd .. you shouldn't be running
		panic(err)
	}
	err = os.Chdir(filepath.Dir(currentDir))
	if err != nil {
		// give up, if you can't do a cd .. you shouldn't be running
		panic(err)
	}
}

func execCmd(printOutput bool, args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	// Extract the command name and arguments
	cmdName := args[0]
	cmdArgs := args[1:]

	// Create the command with the provided arguments
	cmd := exec.Command(cmdName, cmdArgs...)

	// Run the command and capture the combined output
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err := cmd.Run()
	output := stdout.String()
	if printOutput {
		fmt.Println(output)
	}

	if err != nil {
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) > 0 {
			lastLine := lines[len(lines)-1]
			return fmt.Errorf("command execution failed: %v - last line: %s", err, lastLine)
		}
		return fmt.Errorf("command execution failed: %v", err)
	}

	return nil
}

func extractWith7z(file string) error {
	if err := execCmd(false, "7z", "x", file); err == nil {
		return nil
	}
	if err := execCmd(false, "7zz", "x", file); err != nil {
		return err
	}
	return nil
}

func execCmdDirectStdio(args ...string) error {
	if len(args) == 0 {
		return fmt.Errorf("no command provided")
	}

	// Extract the command name and arguments
	cmdName := args[0]
	cmdArgs := args[1:]

	// Create the command with the provided arguments
	cmd := exec.Command(cmdName, cmdArgs...)

	// Set the command's standard input, output, and error to the calling process's standard input, output, and error
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}

	return nil
}

func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func copy(src string, dst string) error {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Ensure the destination directory exists
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0777); err != nil {
		return err
	}

	// Write data to dst
	if err := os.WriteFile(dst, data, 0777); err != nil {
		return err
	}
	return nil
}
