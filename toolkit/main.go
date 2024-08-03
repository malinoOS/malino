package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

var args []string
var Version string = "undefined"
var currentUser *user.User

func main() {
	currentUserTemp, err := user.Current()
	if err != nil {
		fmt.Printf("malino: could not get current username: %v\n", err.Error())
		os.Exit(1)
	}
	currentUser = currentUserTemp

	if stat, err := os.Stat("/home/" + currentUser.Username + "/.malino"); err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir("/home/"+currentUser.Username+"/.malino", 0777)
			if err != nil {
				fmt.Printf("malino: could not create ~/.malino: %v\n", err.Error())
				os.Exit(1)
			}
		} else {
			fmt.Printf("malino: could not check existence of ~/.malino: %v\n", err.Error())
			os.Exit(1)
		}
	} else {
		if !stat.IsDir() {
			fmt.Println("malino: ~/.malino exists, but is a file. Please delete or rename ~/.malino.")
			os.Exit(1)
		}
	}

	args = os.Args[1:]

	if len(args) > 0 { // check if there is any args after "malino"
		switch args[0] { // if there is, then check what it is
		case "help":
			printHelp()
			os.Exit(0)
		case "new":
			if err := newProj(args); err != nil {
				fmt.Printf("malino: error while creating project: %v\n", err.Error())
			}
			os.Exit(0)
		case "build":
			if err := buildProj(); err != nil {
				fmt.Printf("malino: error while building project: %v\n", err.Error())
			}
			os.Exit(0)
		case "run":
			if err := runProj(); err != nil {
				fmt.Printf("malino: error while running project: %v\n", err.Error())
			}
			os.Exit(0)
		case "update-kernel":
			if len(args) > 1 { // check if there is any args after "update-kernel"
				// -no-modules is only here so testing github workflow doesn't take forever.
				// it's also nice not to take bandwith from ubuntu's kernel servers.
				if args[1] == "-no-modules" {
					if err := getKernel(false); err != nil {
						fmt.Printf("malino: error while updating kernel: %v\n", err.Error())
					}
					os.Exit(0)
				}
			}
			if err := getKernel(true); err != nil {
				fmt.Printf("malino: error while updating kernel: %v\n", err.Error())
			}
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
		"malino toolkit v" + Version + "\n\n" +
			"malino help             Shows this help menu\n" +
			"malino new -cs          New malino C# project, uses the name of the folder it's executed in\n" +
			"malino new -go          New malino Go project, uses the name of the folder it's executed in\n" +
			"malino build            Builds your project and creates a .ISO file which can be shared or ran\n" +
			"malino run              Runs your OS in QEMU\n" +
			"malino update-kernel    Downloads the latest Linux kernel. Recommended to do this about every 2 weeks.\n")
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

func goToParentDir() {
	// Get current directory, and CD into the parent directory.
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Chdir(filepath.Dir(currentDir))
	if err != nil {
		panic(err)
	}
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
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stdout

	err := cmd.Run()
	output := stdout.String()
	if printOutput && output != "" {
		fmt.Println(output)
	}

	if err != nil {
		lines := strings.Split(strings.TrimSpace(output), "\n")
		if len(lines) > 0 {
			lastLine := lines[len(lines)-1]
			return fmt.Errorf("%v: command execution failed: %v - last line: %v", strings.Join(args, " "), err, lastLine)
		}
		return fmt.Errorf("command execution failed: %v", err)
	}

	return nil
}

func extractWith7z(file string) error {
	// extract a file with 7z command. if that fails, use the 7zz command. if that fails, return the error.
	if err := execCmd(false, "7z", "x", file); err == nil {
		return nil
	}
	if err := execCmd(false, "7zz", "x", file); err != nil {
		return err
	}
	return nil
}

func downloadFile(url string, filepath string) error {
	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func copyFile(src string, dst string) error {
	// Read all content of src to data, may cause OOM for a large file.
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Ensure the destination directory exists
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, 0777); err != nil {
		return err
	}

	// Write data to dst
	if err := os.WriteFile(dst, data, 0777); err != nil {
		return err
	}
	return nil
}

func copyDirectory(src string, dst string) error {
	// Get file info of the source
	fileInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !fileInfo.IsDir() {
		// If it's a file, call copyFile
		return copyFile(src, dst)
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, fileInfo.Mode()); err != nil {
		return err
	}

	// Read all directory entries
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	// Loop through each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		// Recursively call copyDirectory for each child
		if err := copyDirectory(srcPath, dstPath); err != nil {
			return err
		}
	}

	return nil
}

func combineQuotedStrings(input []string) []string {
	var result []string
	var combined string
	inQuotes := false

	for _, part := range input {
		if strings.HasPrefix(part, "\"") {
			inQuotes = true
			combined = part
		} else if inQuotes {
			combined += " " + part
			if strings.HasSuffix(part, "\"") {
				inQuotes = false
				combined = combined[1 : len(combined)-1] // remove the quotes
				result = append(result, combined)
			}
		} else {
			result = append(result, part)
		}
	}

	// in case the input has a dangling quote, we add the last combined string.
	if inQuotes {
		result = append(result, combined)
	}

	return result
}
