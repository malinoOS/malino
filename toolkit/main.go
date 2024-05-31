package main

import (
	"fmt"
	"os"
)

var args = os.Args[1:]
var Version string = "undefined"

func main() {
	if len(args) > 0 {
		switch args[0] {
		case "help":
			printHelp()
			os.Exit(0)
		case "new":
			newProj(args)
		default:
			fmt.Println("malino: Invalid command")
			printHelp()
			os.Exit(1)
		}
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
