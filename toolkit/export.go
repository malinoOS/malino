package main

import (
	"fmt"
	"os"
)

func exportProj(name string) error {
	fmt.Println("MK iso")
	if err := createAndCD("iso"); err != nil {
		return err
	}
	fmt.Println("MK iso/boot")
	if err := createAndCD("boot"); err != nil {
		os.RemoveAll("iso")
		return err
	}
	fmt.Println("MK iso/boot/grub")
	if err := createAndCD("grub"); err != nil {
		os.RemoveAll("iso")
		return err
	}
	goToParentDir()
	goToParentDir()
	goToParentDir()

	// we really don't want to download vmlinuz like 300 times in one day, kernel.ubuntu.com will probably hate me if i do that
	fmt.Println("CP vmlinuz TO vmlinuz.bak")
	if err := copyFile("vmlinuz", "vmlinuz.bak"); err != nil {
		os.RemoveAll("iso")
		return err
	}
	fmt.Println("MV vmlinuz TO iso/boot/vmlinuz")
	if err := os.Rename("vmlinuz", "iso/boot/vmlinuz"); err != nil {
		os.RemoveAll("iso")
		return err
	}
	fmt.Println("MV vmlinuz.bak TO vmlinuz")
	if err := os.Rename("vmlinuz.bak", "vmlinuz"); err != nil {
		os.RemoveAll("iso")
		return err
	}
	fmt.Println("MV initramfs.cpio.gz TO iso/boot/initramfs.cpio.gz")
	if err := os.Rename("initramfs.cpio.gz", "iso/boot/initramfs.cpio.gz"); err != nil {
		os.RemoveAll("iso")
		return err
	}

	fmt.Println("  W iso/boot/grub/grub.cfg")
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
			"    linux /boot/vmlinuz DOTNET_SYSTEM_GLOBALIZATION_INVARIANT=1\n"+
			"    initrd /boot/initramfs.cpio.gz\n"+
			"}\n"), 0777)
	if err != nil {
		os.RemoveAll("iso")
		return err
	}

	fmt.Println("RUN grub-mkrescue")
	if err := execCmd(true, "grub-mkrescue", "-o", name+".iso", "iso/"); err != nil {
		fmt.Println("RUN /usr/bin/grub-mkrescue")
		if err := execCmd(true, "/usr/bin/grub-mkrescue", "-o", name+".iso", "iso/"); err != nil {
			fmt.Println("RUN /bin/grub-mkrescue")
			if err := execCmd(true, "/bin/grub-mkrescue", "-o", name+".iso", "iso/"); err != nil {
				os.RemoveAll("iso")
				return err
			}
		}
	}

	if err := os.RemoveAll("iso"); err != nil {
		return err
	}

	return nil
}
