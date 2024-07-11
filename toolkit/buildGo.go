package main

import (
	"fmt"

	"github.com/briandowns/spinner"
)

func buildGoProj(spinner *spinner.Spinner, conf []configLine) error {
	fmt.Println(" DL dependencies")
	spinner.Start()
	if err := execCmd(true, "/usr/bin/go", "mod", "tidy"); err != nil {
		spinner.Stop()
		return err
	}
	// make sure we are on the latest version of libmalino
	if err := execCmd(true, "/usr/bin/go", "get", "-u", "github.com/malinoOS/malino/libmalino"); err != nil {
		spinner.Stop()
		return err
	}

	fmt.Println(" GO init")
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
	return nil
}
