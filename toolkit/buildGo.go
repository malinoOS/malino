package main

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

func buildGoProj(spinner *spinner.Spinner) error {
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

	spinner.Stop()
	fmt.Println(" GO init")
	spinner.Start()
	if len(buildflags) == 0 {
		if err := execCmd(true, "/usr/bin/go", "build", "-o", "mInit", "-ldflags", "-X main.Version="+time.Now().Format("060102")); err != nil {
			spinner.Stop()
			return err
		}
	} else {
		if err := execCmd(true, append([]string{"/usr/bin/go", "build", "-o", "mInit", "-ldflags", "-X main.Version=" + time.Now().Format("060102")}, buildflags...)...); err != nil {
			spinner.Stop()
			return err
		}
	}
	spinner.Stop()
	return nil
}
