package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type configLine struct {
	hasAnything bool
	err         error
	operation   string
	args        []string
}

func parseConfigLine(line string, lineNum int) configLine {
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
			return configLine{false, fmt.Errorf("malino.cfg: line %v: line does not contain 3 words, which is required for include operation", lineNum), "", nil}
		}
		op = "include"

	case "buildflags":
		return configLine{true, nil, "buildflags", combineQuotedStrings(words[1:])}

	case "verfmt":
		// verfmt is deprecated now, does nothing
		// it will always be yymmdd, as i've never seen anybody use any other option, or not use it even.
		return configLine{false, nil, "", nil}

	case "lang":
		if len(words) != 2 {
			return configLine{false, fmt.Errorf("malino.cfg: line %v: line does not contain 2 words, which is required for lang operation", lineNum), "", nil}
		}
		op = "lang"

	default:
		return configLine{false, fmt.Errorf("malino.cfg: line %v: invalid operation", lineNum), "", nil}
	}

	return configLine{true, nil, op, words[1:]}
}

func handleIncludeLine(line configLine) error {
	if !line.hasAnything {
		// It really isn't.
		// If you encounter this error, it's because the parseConfigLine() function said that the line doesn't have anything,
		// But the build.go file that first parses the config, ignored that, and decided to still call a handler.
		// This should never happen, but sometimes it does.
		// This is fine, just check build.go, and see if it's checking line.hasAnything correctly.
		return fmt.Errorf("the entire configuration parser is broken. good luck")
	}

	if line.operation == "include" {
		curDir := "undefined"
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			curDir = dir
		}

		line.args[0] = strings.Replace(line.args[0], "./", filepath.Dir(curDir)+"/", 1)
		fmt.Printf("INC %v AS %v\n", line.args[0], line.args[1])

		if strings.HasPrefix(line.args[0], "https://") {
			if err := downloadFile(line.args[0], curDir+line.args[1]); err != nil {
				return err
			}
			return nil
		} else if strings.HasPrefix(line.args[0], "dir///") {
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
