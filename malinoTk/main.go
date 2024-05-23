package main

import (
	"fmt"
	"os"
)

var args = os.Args[1:]

func main() {
	if len(args) > 0 {
		if args[0] == "help" {
			printHelp()
			os.Exit(0)
		} else if args[0] == "new" {
			if len(args) > 1 {
				err := newProj(args[1])
				if err != nil {
					fmt.Printf("err: %v", err.Error())
					os.Exit(1)
				} else {
					wd, err := os.Getwd()
					if err != nil {
						fmt.Printf("Project \"%v\" created.", args[1])
					} else {
						fmt.Printf("Project \"%v/%v\" created.", wd, args[1])
					}
				}
			} else {
				printHelp()
				os.Exit(0)
			}
		} else if args[0] == "build" {
			err := buildProj()
			if err != nil {
				fmt.Printf("err: %v", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("build successful")
			}
		} else if args[0] == "run" {
			err := runProj()
			if err != nil {
				fmt.Printf("err: %v", err.Error())
				os.Exit(1)
			} else {
				fmt.Printf("run successful")
			}
		}
	} else {
		printHelp()
		os.Exit(0)
	}
}

func printHelp() {
	fmt.Print(
		"malino toolkit\n\n" +
			"malino help = shows this help menu\n" +
			"malino new [name] = new project, creates folder and go module with name [name]\n" +
			"malino build = builds a qcow2 disk image of your OS\n" +
			"malino run = runs your built qcow2 disk image in QEMU\n" +
			"malino export = exports your OS into a .ISO file which can be shared or burned onto a CD\n")
}
