<img src="https://winksplorer.net/malino.png" width="250">

[discord server](https://discord.gg/2yfxxfNT6F)

# malino
The malino Linux-based OS development toolkit. In other words, you get to create an OS based on the Linux kernel.

In other-other words, you get to write an initramfs busybox-replacement in Go, and we take care of the annoying parts. This is your toolchain now.

# How to install
> malino is for linux only.

## From the debian package
1. Download the .deb file from the releases tab.
2. Run `sudo dpkg -i ~/Downloads/malino-stable-*`.

## From source
1. Download requirements. On Debian-based distros, run `sudo apt install golang-go qemu-system-x86 qemu-utils p7zip`.
2. Clone the repository with `git clone https://github.com/malinoOS/malino`
3. Cd into the repo and run `make`.
4. Now malino is in /usr/bin, and you can use it.
5. To create a new malino project, run `malino new <your project name here>`, replacing `<your project name here>` with the name of your project.
6. To compile your project, run `malino build`.
7. To run your project in QEMU, run `malino run`.
8. To export your project into a .ISO file for you to run on real hardware and share, run `malino export`. If you want UEFI support, run `malino export -efi`.

## libmalino
libmalino is the Go module that your OS imports, so you don't need 50 lines just to read a line from the user.

## malino
malino is the toolkit and command you use to create projects, build, export, etc...
