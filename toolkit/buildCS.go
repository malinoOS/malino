package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/briandowns/spinner"
)

func buildCSProj(spinner *spinner.Spinner, conf []configLine) error {
	fmt.Println(" C# init")
	spinner.Start()
	name := "undefined"
	curDir := "undefined"
	if dir, err := os.Getwd(); err != nil {

		return err
	} else {
		name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
		curDir = dir
	}

	currentUser, err := user.Current()
	if err != nil {
		spinner.Stop()
		return err
	}

	if _, err := os.Stat("/home/" + currentUser.Username + "/.bflat/bflat"); os.IsNotExist(err) {
		if err := DownloadCSCompiler("/home/" + currentUser.Username); err != nil {
			return err
		}
	}

	buildFlagsExist := false
	for _, line := range conf {
		if line.operation == "buildflags" {
			buildFlagsExist = true
			// .net moment
			csFiles, err := getCsFiles(curDir)
			if err != nil {
				return err
			}
			buildCmd := append([]string{"/home/" + currentUser.Username + "/.bflat/bflat", "build", "-o", "mInit", "-r/opt/malino/libmalino-cs.dll"}, csFiles...)
			if err := execCmd(true, append(buildCmd, line.args...)...); err != nil {
				spinner.Stop()
				return err
			}
			break
		} else if line.operation == "verfmt" {
			ver, err := handleVerfmtLine(line)
			if err != nil {
				return err
			}
			err = os.WriteFile("malino.generated.cs", []byte(
				"namespace "+name+" {\n"+
					"	// Don't you dare mess with this.\n"+
					"	public class MalinoAutoGenerated {\n"+
					"		public static string OSVersion = \""+ver+"\";\n"+
					"	}\n"+
					"}"), 0777)
			if err != nil {
				spinner.Stop()
				return err
			}
			break
		}
	}
	if !buildFlagsExist {
		csFiles, err := getCsFiles(curDir)
		if err != nil {
			return err
		}
		buildCmd := append([]string{"/home/" + currentUser.Username + "/.bflat/bflat", "build", "-o", "mInit", "-r/opt/malino/libmalino-cs.dll"}, csFiles...)
		if err := execCmd(true, buildCmd...); err != nil {
			spinner.Stop()
			return err
		}
	}

	spinner.Stop()
	return nil
}

func DownloadCSCompiler(homeDirectory string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := downloadFile("https://github.com/bflattened/bflat/releases/download/v8.0.2/bflat-8.0.2-linux-glibc-x64.tar.gz", "bflat.tar.gz"); err != nil {
		return err
	}

	if err := createAndCD(homeDirectory + "/.bflat"); err != nil {
		return err
	}

	if err := execCmd(true, "tar", "-xzf", wd+"/bflat.tar.gz"); err != nil {
		return err
	}

	if err := os.Chdir(wd); err != nil {
		return err
	}

	return nil
}

func getCsFiles(dir string) ([]string, error) {
	var csFiles []string

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".cs" {
			csFiles = append(csFiles, file.Name())
		}
	}

	return csFiles, nil
}
