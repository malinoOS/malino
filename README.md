<p align="center"><img src="https://winksplorer.net/malino-whitetext.png" width="500"></p>

<p align="center">
<a href="https://discord.gg/2yfxxfNT6F">
    <img src="https://img.shields.io/badge/discord-server-blue?style=for-the-badge&logo=discord">
</a>
<img src="https://img.shields.io/github/languages/code-size/malinoOS/malino?style=for-the-badge&logo=files"/>
<a href="https://github.com/malinoOS/malino/blob/master/LICENSE">
    <img src="https://img.shields.io/github/license/malinoOS/malino?style=for-the-badge&logo=gnu"/>
</a>
</p>

## How to install
### [GitHub Wiki: Installation](https://github.com/malinoOS/malino/wiki/Installation)

## How to use
### [GitHub Wiki: Toolkit usage](https://github.com/malinoOS/malino/wiki/Toolkit-usage)

# Directory structure

## libmalino
libmalino is the Go module that your OS imports, so you don't need 50 lines just to read a line from the user.

Include it in your Go file with `import "github.com/malinoOS/malino/libmalino"`.

## libmalino-cs
libmalino-cs is libmalino but for C#. It uses .NET 8.0 to compile, and is placed in `/opt/malino/libmalino-cs.dll`.

malino automatically "links" libmalino-cs with your project if you have your project configured to build for C#.

## libmsb
MSB stands for "Malino Syscall Bridge". This is only used with C# projects, and it's used to allow C# to make Linux system calls, since for some reason it can't by default. And it uses `clang` to build since this is a syscall bridge and must be written in C.

## malino
malino is the toolkit and command you use to create projects, build, export, etc...
