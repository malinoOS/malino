package main

import (
	"fmt"
	"os"
)

var args = os.Args[1:]
var Version string = "undefined"

func main() {
	if len(args) > 0 {
		if args[0] == "help" {
			printHelp()
			os.Exit(0)
		} else if args[0] == "new" {
			if len(args) > 1 {
				if len(args) > 2 {
					err := newProjHere(args[1])
					if err != nil {
						fmt.Printf("err: %v\n", err.Error())
						os.Exit(1)
					} else {
						fmt.Printf("Project \"%v\" created.\n", args[1])
					}
				} else {
					err := newProj(args[1])
					if err != nil {
						fmt.Printf("err: %v\n", err.Error())
						os.Exit(1)
					} else {
						wd, err := os.Getwd()
						if err != nil {
							fmt.Printf("Project \"%v\" created.\n", args[1])
						} else {
							fmt.Printf("Project \"%v/%v\" created\n", wd, args[1])
						}
					}
				}
			} else {
				printHelp()
				os.Exit(0)
			}
		} else if args[0] == "build" {
			err := buildProj()
			if err != nil {
				fmt.Printf("err: %v\n", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("build successful\n")
			}
		} else if args[0] == "run" {
			err := runProj()
			if err != nil {
				fmt.Printf("err: %v\n", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("run successful\n")
			}
		}
	} else {
		printHelp()
		os.Exit(0)
	}
}

func printHelp() {
	fmt.Print(
		"malino toolkit v" + Version + "\n\n" +
			"malino help                Shows this help menu\n" +
			"malino new [name] [here]   New project, creates folder and go module with name [name], type \"here\" for [here] to not create a new folder.\n" +
			"malino build               Builds a qcow2 disk image of your OS\n" +
			"malino run                 Runs your built qcow2 disk image in QEMU\n" +
			"malino export              Exports your OS into a .ISO file which can be shared or burned onto a CD\n")
}
