package main

import (
	"fmt"
	"os"
	"os/exec"
)

func buildProj() error {
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory does not contain a valid malino project")
	}

	cmd := exec.Command("/usr/bin/go", "build")
	stdout, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Println(string(stdout))

	return nil
}
