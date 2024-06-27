<p align="center"><img src="https://github.com/malinoOS/malino/assets/49623720/c764dd50-c0cd-4440-993a-49373ebba912" width="400"></p>

<p align="center">
<a href="https://discord.gg/2yfxxfNT6F"><img src="https://img.shields.io/badge/chat-on_discord-blue?style=for-the-badge&logo=discord"></a>
<img src="https://img.shields.io/github/languages/code-size/malinoOS/malino?style=for-the-badge&logo=files"/>
<a href="https://github.com/malinoOS/malino/releases/latest"><img src="https://img.shields.io/github/v/release/malinoOS/malino?style=for-the-badge&logo=debian" /></a>
<a href="https://github.com/malinoOS/malino/blob/master/LICENSE"><img src="https://img.shields.io/github/license/malinoOS/malino?style=for-the-badge&logo=mozilla"/></a>
</p>

### malino is a toolkit that allows people to create their own operating systems, easily.

#### It supports both Go & C#, and you get to use Linux as your base.

#### And also has a library that helps you make an OS with the toolkit.

##### (in beta)

## Features

- Direct Linux system call access
- Advanced file system, supports many filesystems, and works on real hardware
- Most features found in both the C# and Go standard library
- BIOS & EFI support on real hardware, almost all features work on real hardware
- Framebuffer support to the point where it can [run DOOM](https://youtu.be/JERv-ocRCW4).
- Including files in the system, allows lots of apps (with their libraries) to be ran (including apps like `ffmpeg`)
- Faster than Cosmos in almost every way


### How to install
#### [GitHub Wiki: Installation](https://github.com/malinoOS/malino/wiki/Installation)

### How to use
#### [GitHub Wiki: Getting Started](https://github.com/malinoOS/malino/wiki/Getting-Started)
#### [GitHub Wiki: Toolkit usage](https://github.com/malinoOS/malino/wiki/Toolkit-usage)

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
