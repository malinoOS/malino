package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
)

func exportProj() error {
	// Initialize the spinner (loading thing).
	spinner := spinner.New(spinner.CharSets[14], 100*time.Millisecond)

	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("current directory doesn't contain a project")
	}

	name := "undefined"
	if dir, err := os.Getwd(); err != nil {
		return err
	} else {
		name = strings.Split(dir, "/")[len(strings.Split(dir, "/"))-1] // "name = split by / [len(split by /) - 1]" basically.
	}

	if _, err := os.Stat("vmlinuz"); os.IsNotExist(err) {
		if err := buildProj(); err != nil {
			return err
		}
	}

	if _, err := os.Stat("initramfs.cpio.gz"); os.IsNotExist(err) {
		if err := buildProj(); err != nil {
			return err
		}
	}

	spinner.Start()
	fmt.Println("MK iso")
	if err := createAndCD("iso"); err != nil {
		spinner.Stop()
		return err
	}
	fmt.Println("MK iso/boot")
	if err := createAndCD("boot"); err != nil {
		spinner.Stop()
		return err
	}
	fmt.Println("MK iso/boot/grub")
	if err := createAndCD("grub"); err != nil {
		spinner.Stop()
		return err
	}
	goToParentDir()
	goToParentDir()
	goToParentDir()
	spinner.Stop()

	spinner.Start()
	// we really don't want to download vmlinuz like 300 times in one day, kernel.ubuntu.com will probably hate me if i do that
	fmt.Println("CP vmlinuz TO vmlinuz.bak")
	if err := copyFile("vmlinuz", "vmlinuz.bak"); err != nil {
		spinner.Stop()
		return err
	}
	fmt.Println("MV vmlinuz TO iso/boot/vmlinuz")
	if err := os.Rename("vmlinuz", "iso/boot/vmlinuz"); err != nil {
		spinner.Stop()
		return err
	}
	fmt.Println("MV vmlinuz.bak TO vmlinuz")
	if err := os.Rename("vmlinuz.bak", "vmlinuz"); err != nil {
		spinner.Stop()
		return err
	}
	fmt.Println("MV initramfs.cpio.gz TO iso/boot/initramfs.cpio.gz")
	if err := os.Rename("initramfs.cpio.gz", "iso/boot/initramfs.cpio.gz"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	fmt.Println("  W iso/boot/grub/grub.cfg")
	spinner.Start()
	err := os.WriteFile("iso/boot/grub/grub.cfg", []byte(
		"set default=0\n"+
			"set timeout=0\n\n"+
			"if [ \"${grub_platform}\" = \"efi\" ]; then\n"+
			"    insmod efi_gop\n"+
			"    insmod efi_uga\n"+
			"else\n"+
			"    insmod vbe\n"+
			"fi\n"+
			"insmod font\n"+
			"if loadfont /boot/grub/fonts/unicode.pf2; then\n"+
			"    insmod gfxterm\n"+
			"    set gfxmode=auto\n"+
			"    set gfxpayload=keep\n"+
			"    terminal_output gfxterm\n"+
			"fi\n\n"+
			"menuentry '"+name+"' --class os {\n"+
			"    insmod gzio\n"+
			"    insmod part_msdos\n"+
			"    linux /boot/vmlinuz\n"+
			"    initrd /boot/initramfs.cpio.gz\n"+
			"}\n"), 0777)
	if err != nil {
		spinner.Stop()
		return err
	}

	fmt.Println("RUN grub-mkrescue")
	if err := execCmd(false, "grub-mkrescue", "-o", name+".iso", "iso/"); err != nil {
		spinner.Stop()
		return err
	}

	if err := os.RemoveAll("iso"); err != nil {
		spinner.Stop()
		return err
	}
	spinner.Stop()

	return nil
}
