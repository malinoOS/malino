package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type configLine struct {
	err       error
	operation string
	args      []string
}

func parseConfigLine(line string, lineNum int) configLine {
	// check if line is empty or is a comment, if it is, just say it didn't do anything
	if line == "" || strings.HasPrefix(line, "#") {
		return configLine{nil, "", nil}
	}
	// split by spaces, throw error if there is not 3 words since that's the syntax
	words := strings.Split(line, " ")

	op := ""

	switch words[0] {
	case "include":
		if len(words) != 3 {
			return configLine{fmt.Errorf("malino.cfg: line %v: line does not contain 3 words, which is required for include operation", lineNum), "", nil}
		}
		op = "include"

	case "modpack":
		if len(words) != 2 {
			return configLine{fmt.Errorf("malino.cfg: line %v: line does not contain 2 words, which is required for modpack operation", lineNum), "", nil}
		}
		op = "modpack"

	case "buildflags":
		return configLine{nil, "buildflags", combineQuotedStrings(words[1:])}

	case "verfmt":
		// verfmt is deprecated now, does nothing
		// it will always be yymmdd, as i've never seen anybody use any other option, or not use it even.
		return configLine{nil, "", nil}

	case "lang":
		if len(words) != 2 {
			return configLine{fmt.Errorf("malino.cfg: line %v: line does not contain 2 words, which is required for lang operation", lineNum), "", nil}
		}
		op = "lang"

	default:
		return configLine{fmt.Errorf("malino.cfg: line %v: invalid operation", lineNum), "", nil}
	}

	return configLine{nil, op, words[1:]}
}

func handleIncludeLine(line configLine) error {
	if line.operation == "include" {
		curDir := "undefined"
		if dir, err := os.Getwd(); err != nil {
			return err
		} else {
			curDir = dir
		}

		line.args[0] = strings.Replace(line.args[0], "./", filepath.Dir(curDir)+"/", 1)
		fmt.Printf("INC %v AS %v\n", line.args[0], line.args[1])

		if strings.HasPrefix(line.args[0], "https://") { // if https link
			if err := downloadFile(line.args[0], curDir+line.args[1]); err != nil {
				return err
			}
			return nil
		} else if strings.HasPrefix(line.args[0], "dir///") { // dir/// means include an entire directory
			if err := copyDirectory(line.args[0][6:], curDir+line.args[1]); err != nil {
				if !os.IsNotExist(err) {
					return err
				} else {
					fmt.Printf("WARNING: included directory %s doesn't exist on your host system!\n", line.args[0][6:])
				}
			}
			return nil
		}

		if err := copyFile(line.args[0], curDir+line.args[1]); err != nil { // otherwise, just copy file
			if !os.IsNotExist(err) {
				return err
			} else {
				fmt.Printf("WARNING: included file %s doesn't exist on your host system!\n", line.args[0])
			}
		}

	} else {
		return fmt.Errorf("include handler called for non-include operation")
	}

	return nil
}

func handleModpackLine(line configLine) error {
	if line.operation == "modpack" {
		// init module table
		configLookup := map[string][]configLine{
			"usb": {
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/usb/host/xhci-plat-hcd.ko", currentUser.Username), "/modules/xhci-plat-hcd.ko"}},
			},
			"usbhid": {
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/hid/usbhid/usbkbd.ko", currentUser.Username), "/modules/usbkbd.ko"}},
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/hid/usbhid/usbmouse.ko", currentUser.Username), "/modules/mouse/usbmouse.ko"}},
			},
			"mouse": {
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/input/mouse/gpio_mouse.ko", currentUser.Username), "/modules/mouse/gpio_mouse.ko"}},
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/input/mouse/psmouse.ko", currentUser.Username), "/modules/mouse/psmouse.ko"}},
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/input/mouse/sermouse.ko", currentUser.Username), "/modules/mouse/sermouse.ko"}},
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/input/mouse/synaptics_i2c.ko", currentUser.Username), "/modules/mouse/synaptics_i2c.ko"}},
				{nil, "include", []string{fmt.Sprintf("/home/%s/.malino/modules/drivers/input/mouse/synaptics_usb.ko", currentUser.Username), "/modules/mouse/synaptics_usb.ko"}},
			},
		}

		configLines, exists := configLookup[line.args[0]]
		if !exists {
			return fmt.Errorf("module pack \"%s\" doesn't exist", line.args[0])
		}

		for _, line := range configLines {
			// "That's right, we're gonna cheat."
			// A modpack is basically a collection of includes.
			if err := handleIncludeLine(line); err != nil {
				return err
			}
		}

		return nil
	} else {
		return fmt.Errorf("modpack handler called for non-modpack operation")
	}
}
