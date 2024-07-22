using System;
using System.Reflection.Metadata;
using System.Runtime.InteropServices;
using System.Text;

namespace libmalino;

// please shut up
#pragma warning disable CA1401
#pragma warning disable CA2101
#pragma warning disable SYSLIB1054

/// <summary>
/// These are bindings for libmsb, a core library that malino needs for the C# variant of libmalino.
///
/// Since C# cannot do syscalls, we have to make glue code in C that will make the syscalls, and make bindings for it.
/// libmsb does exactly this. While malino could just make bindings for libc's syscalls, P/Invoking libc is probably not a good idea.
/// </summary>
public class MsbBindings
{
    /// <summary>
    /// A bind for sync() in libc: https://man7.org/linux/man-pages/man2/sync.2.html.
    ///
    /// Does not return anything.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_sync")]
    public static extern void Sync();

    /// <summary>
    /// A bind for linux's reboot system call: https://man7.org/linux/man-pages/man2/reboot.2.html.
    ///
    /// If the command stops or restarts the system, then a successful call should never return.
    /// For other commands, 0 is for success, and if it fails, it will return the errno value.
    ///
    /// Use `libmalino.Errno.Msg[errorCode]` to see what it's complaining about.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_reboot")]
    public static extern int Reboot(uint cmd);

    /// <summary>
    /// A bind for write() in libc: https://man7.org/linux/man-pages/man2/write.2.html.
    ///
    /// If successful, it returns how many bytes it wrote (Read about success values using `man write.2`).
    /// If it failed, then it will return the negative value of it's errno, for example -42 would be 42 ENOMSG.
    ///
    /// Use `libmalino.Errno.Msg[-errorCode]` to see what it's complaining about.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_write")]
    public static extern long Write(uint fd, string buf);

    [DllImport("libmsb", EntryPoint = "msb_read")]
    private static extern long ReadDirect(uint fd, StringBuilder buf, ulong count);

    /// <summary>
    /// A bind for read() in libc: https://man7.org/linux/man-pages/man2/read.2.html.
    ///
    /// If successful, it returns how many bytes it read (Read about success values using `man read.2`).
    /// If it failed, then it will return the negative value of it's errno, for example -42 would be 42 ENOMSG.
    ///
    /// Use `libmalino.Errno.Msg[-errorCode]` to see what it's complaining about.
    /// </summary>
    public static string Read(uint fd, ulong count) {
        StringBuilder sb = new();
        ReadDirect(fd, sb, count);
        return sb.ToString();
    }

    /// <summary>
    /// A bind for mount() in libc: https://man7.org/linux/man-pages/man2/mount.2.html.
    ///
    /// 0 is for success, and if it fails, it will return the errno value.
    /// Use `libmalino.Errno.Msg[errorCode]` to see what it's complaining about.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_mount")]
    public static extern int Mount(string source, string target, string fstype, ulong flags, string data);

    /// <summary>
    /// A bind for umount2() in libc: https://man7.org/linux/man-pages/man2/umount2.2.html.
    ///
    /// 0 is for success, and if it fails, it will return the errno value.
    /// Use `libmalino.Errno.Msg[errorCode]` to see what it's complaining about.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_umount2")]
    public static extern int Unmount(string target, int flags);

    /// <summary>
    /// A bind for getpid() in libc: https://man7.org/linux/man-pages/man2/getpid.2.html.
    ///
    /// This should never fail.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_getpid")]
    public static extern int GetPID();

    /// <summary>
    /// A combination of fork() (https://man7.org/linux/man-pages/man2/fork.2.html) and execve() (https://man7.org/linux/man-pages/man2/execve.2.html) in libc.
    ///
    /// If successful AND wait=true, it will return the status code of the application after it has exited.
    /// IF successful AND wait=false, it will return 0.
    /// If it failed, then it will return the negative value of it's errno, for example -42 would be 42 ENOMSG.
    ///
    /// Use `libmalino.Errno.Msg[-errorCode]` to see what it's complaining about.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_forkexec")]
    public static extern int ForkExec(string path, string[] argv, string[] envp, bool wait);

    /// <summary>
    /// Calls a direct Linux system call.
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_dsc")]
    public static extern int Syscall(long rax, long rdi, long rsi, long rdx, long r10, long r8, long r9);

    /// <summary>
    /// Loads a Linux kernel module (.ko format). USE LIBMALINO.MALINO.LOADKERNELMODULE() INSTEAD!!!!
    /// </summary>
    [DllImport("libmsb", EntryPoint = "msb_loadko")]
    public static extern int LoadKernelModule(string path, string param);
}

// please unshut up
#pragma warning restore CA1401
#pragma warning restore CA2101
#pragma warning restore SYSLIB1054