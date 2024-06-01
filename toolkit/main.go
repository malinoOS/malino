package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
			newProj(args)
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
		"malino toolkit (rewrite branch) v" + Version + "\n\n" +
			"malino help         Shows this help menu\n" +
			"malino new [name]   New project, creates folder and go module with name [name]\n" +
			"malino new          New project, does not make folder, and uses the name of the folder it's executed in\n" +
			"malino build        Builds a disk image of your OS\n" +
			"malino run          Runs your built  disk image in QEMU\n" +
			"malino export       Exports your OS into a .ISO file which can be shared or burned onto a CD\n")
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

func goToParentDir() error {
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
	return nil
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
	stdout, err := cmd.CombinedOutput()
	if printOutput {
		fmt.Println(string(stdout))
	}
	if err != nil {
		return fmt.Errorf("command execution failed: %v", err)
	}
	return nil
}
