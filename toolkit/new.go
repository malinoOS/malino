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
			return err
		}
	}
	spinner.Stop()

	fmt.Println("Creating Go project...")
	spinner.Start()
	if err := execCmd(false, "/usr/bin/go", "mod", "init", name); err != nil { // init the go module
		return err
	}
	err := os.WriteFile("main.go", []byte(
		"package main\n\n"+
			"import (\n"+
			"	\"fmt\"\n"+
			")\n\n"+
			"func main() {\n"+
			"	fmt.Println(\"malino (project "+name+") booted successfully. Type a line of text to get it echoed back.\")\n"+
			/*"	for {\n"+
			"		fmt.Print(\"Input: \")\n"+
			"		input := libmalino.UserLine()\n"+
			"		fmt.Println(\"Text typed: \" + input)\n"+
			"	}\n"+*/
			"}"), 0777)
	if err != nil {
		return err
	}
	spinner.Stop()

	fmt.Println("Creating VM...")
	spinner.Start()
	// Doing this in a Makefile is easier. deal with it.
	err = os.WriteFile("Makefile", []byte(
		"SHELL := /bin/bash\n\n"+
			"all:\n"+
			"ifeq ($(MALINO_KEY), thisKeyIsForMakingSureThatAllVariablesThatTheMalinoEnvironmentUsesIsPresentSoItDoesNotRunIfYouJustTypeMakeHere)\n"+
			"	sudo modprobe nbd max_part=8\n"+
			"	qemu-img create -f qcow2 /home/$(shell whoami)/.malinoVM_$(MALINO_PROJ).qcow2\n"+
			"	sudo qemu-nbd -c /dev/nbd3 /home/$(shell whoami)/.malinoVM_$(MALINO_PROJ).qcow2\n"+
			"	echo -e \"o\\nn\\np\\n\\n\\n\\nw\" | sudo fdisk /dev/nbd3\n"+ // readable code
			"	sudo mkfs -t ext4 /dev/nbd3p1\n"+
			"	mkdir disk\n"+
			"	sudo mount -t ext4 /dev/nbd3p1 disk\n"+
			"	sudo mkdir -pv disk/{bin,sbin,etc,lib,lib64,var,dev,proc,sys,run,tmp,boot}\n"+
			"	sudo mknod -m 600 disk/dev/console c 5 1\n"+
			"	sudo mknod -m 600 disk/dev/tty c 5 1\n"+
			"	sudo mknod -m 666 disk/dev/null c 1 3\n"+
			"	sudo cp $$(ls -t /boot/vmlinuz* | head -n1) disk/boot/\n"+ // hmm
			"	sudo cp $$(ls -t /boot/initrd* | head -n1) disk/boot/\n"+ // hmmmm
			"	sudo grub-install /dev/nbd3 --skip-fs-probe --boot-directory=disk/boot --target=i386-pc\n"+
			"	printf \"set default=0\\nset timeout=1\\n\\nmenuentry \\\"golinux\\\" {\\n    linux $$(ls -t /boot/vmlinuz* | head -n1) root=/dev/sda1 ro\\n    initrd $$(ls -t /boot/initrd* | head -n1)\\n}\" | sudo tee disk/boot/grub/grub.cfg\n"+
			"	printf \"[init]\\nprintSplashMessage = true\\nremountRootPartitionAsWritable = true\\nmalinoMode = true\\nexec = /bin/fallsh\" | sudo tee disk/etc/init.ini\n"+
			"	sudo umount disk\n"+
			"	rm -rf disk\n"+
			"	qemu-nbd -d /dev/nbd3\n"+
			"else\n"+
			"	$(error This environment isn't the malino builder)\n"+
			"endif\n"), 0777)
	if err != nil {
		return err
	}

	if hasNameArg {
		goToParentDir()
	}

	return nil
}
