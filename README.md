<img src="https://winksplorer.net/malino.png" width="250">

# malino
The Malino Linux-based OS development toolkit

# How to use
> malino is for linux only.

1. Clone the repository with `git clone https://github.com/malinoOS/malino`
2. Cd into the repo and run `make`.
3. Now malino is in /usr/bin, and you can use it.
4. To create a new malino project, run `malino new <your project name here>`, replacing `<your project name here>` with the name of your project.
5. To compile your project, run `malino build`.
6. To run your project in QEMU, run `malino run`.

`malino publish` is not implemented yet.

## libmalino
libmalino is the Go module that your OS imports, so you don't need 50 lines just to read a line from the user.

## malino
malino is the toolkit and command you use to create projects, build, export, etc...