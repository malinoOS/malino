package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func newProj(args []string) error {
	// Initialize the spinner (loading thing).
	spinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	// Find the name that the project should have. It checks if there is anything after "new" in "malino new".
	// "malino new" = the project name will be the name of the current folder.
	// "malino new test" = the project will have the name of "test".
	name := "undefined"
	hasNameArg := len(args) == 2
	if hasNameArg {
		name = args[1]
	} else {
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
		}
	}
	fmt.Println("Found name: " + name)

	fmt.Println("Creating directories...")
	spinner.Start()
	if hasNameArg { // root directory for project. only create if name is specified in args.
		if err := createAndCD(name); err != nil {
			spinner.Stop()
			return err
		}
	}
	spinner.Stop()

	fmt.Println("Creating Go project...")
	spinner.Start()
	if err := execCmd(false, "/usr/bin/go", "mod", "init", name); err != nil { // init the go module
		spinner.Stop()
		return err
	}
	err := os.WriteFile("main.go", []byte(
		"package main\n\n"+
			"import (\n"+
			"	\"github.com/malinoOS/malino/libmalino\"\n"+
			"	\"fmt\"\n"+
			")\n\n"+
			"func main() {\n"+
			"	defer libmalino.resetTerminalMode()\n"+
			"	fmt.Println(\"malino (project "+name+") booted successfully. Type a line of text to get it echoed back.\")\n"+
			"	for { // Word of advice: Never let this app exit. Always end in an infinite loop or shutdown.\n"+
			"		fmt.Print(\"Input: \")\n"+
			"		input := libmalino.UserLine()\n"+
			"		fmt.Println(\"Text typed: \" + input)\n"+
			"	}\n"+
			"}"), 0777)
	if err != nil {
		spinner.Stop()
		return err
	}

	err = os.WriteFile(".gitignore", []byte("vmlinuz\ninitramfs.cpio.gz\n"+name+".iso"), 0777)
	if err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	if hasNameArg {
		goToParentDir()
	}

	fmt.Println("Done.")

	return nil
}
