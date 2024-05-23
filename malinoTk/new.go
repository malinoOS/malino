package main

import (
	"archive/zip"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func newProj(name string) error {
	println("creating directory...")
	err := os.Mkdir(name, 0777)
	if err != nil {
		return err
	}

	err = os.Chdir(name)
	if err != nil {
		return err
	}

	println("creating Go project...")
	cmd := exec.Command("/usr/bin/go", "mod", "init", name)
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	println("creating template code...")
	//err = os.WriteFile(fmt.Sprintf("%v/config.malino", name), fmt.Sprintf("[proj]\nname = %v\n")) // nerds say i should "use + instead of sprintf!!!" i literally don't care
	err = os.WriteFile("main.go", []byte(
		"package main\n\n"+
			"import (\n"+
			"	\"github.com/malinoOS/malino/libmalino\"\n"+
			"	\"fmt\"\n"+
			")\n\n"+
			"func main() {\n"+
			"	fmt.Println(\"welcome to "+name+"!\")\n"+
			"	libmalino.test()\n"+
			"	for {}\n"+
			"}"), 0777)
	if err != nil {
		return err
	}

	println("downloading golinux...")
	err = DownloadFile("https://github.com/malinoOS/golinux/archive/refs/heads/main.zip", "golinux.zip")
	if err != nil {
		return err
	}

	println("unzipping golinux...")
	err = Unzip("golinux.zip", ".") // uhhh
	if err != nil {
		return err
	}

	println("deleting zip...")
	err = os.Remove("golinux.zip")
	if err != nil {
		return err
	}

	err = os.Chdir("golinux-main")
	if err != nil {
		return err
	}

	println("setting up golinux (this will take a while)...")
	cmd = exec.Command("/usr/bin/make", "createVM", "clean", "prepare", "buildInit", "buildFallsh", "install")
	_, err = cmd.Output()
	if err != nil {
		return err
	}

	// remember to go back to the root where malino command was executed!
	goToParentFolder()
	goToParentFolder()

	return nil
}

func DownloadFile(url string, filepath string) error {
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

func Unzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			return err
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	return nil
}

func goToParentFolder() error {
	currentDir, err := os.Getwd()
	if err != nil {
		// give up, if you can't do a cd .. you shouldn't be running
		println(err.Error())
		os.Exit(1)
	}
	err = os.Chdir(filepath.Dir(currentDir))
	if err != nil {
		// give up, if you can't do a cd .. you shouldn't be running
		println(err.Error())
		os.Exit(1)
	}
	return nil
}
