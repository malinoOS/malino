using System;
using System.IO;

namespace libmalino;

/// <summary>
/// Core functions.
/// </summary>
public class malino {
    /// <summary>
    /// Syncs filesystems then shuts down the computer. If it fails then it will throw an exception of the Errno message.
    /// </summary>
    public static void ShutdownComputer() {
        Console.WriteLine("syncing disks...");
        MsbBindings.Sync();
        Console.WriteLine("shutting down...");
        int val = MsbBindings.Reboot((uint)LINUX_REBOOT.CMD_POWER_OFF);
        if (val != 0)
            throw new Exception(Errno.GetStringErr(val));
    }

    /// <summary>
    /// Syncs filesystems then reboots the computer. If it fails then it will throw an exception of the Errno message.
    /// </summary>
    public static void RebootComputer() {
        Console.WriteLine("syncing disks...");
        MsbBindings.Sync();
        Console.WriteLine("shutting down...");
        int val = MsbBindings.Reboot((uint)LINUX_REBOOT.CMD_RESTART);
        if (val != 0)
            throw new Exception(Errno.GetStringErr(val));
    }

    /// <summary>
    /// Returns the amount of seconds since Linux has booted.
    /// </summary>
    public static int SystemUptimeAsInt() {
        string dat = File.ReadAllText("/proc/uptime");
        string[] stringRep = dat.Split(" ")[0].Split(".");
        return int.Parse(stringRep[0]);
    }

    /// <summary>
    /// Returns the amount of seconds to the hundredths since Linux has booted.
    /// </summary>
    public static float SystemUptimeAsFloat() {
        string dat = File.ReadAllText("/proc/uptime");
        string stringRep = dat.Split(" ")[0];
        return float.Parse(stringRep);
    }

    /// <summary>
    /// chdirs into a directory, and spawns a process there, without killing the parent process.
    /// </summary>
    public static int SpawnProcess(string path, string startDir, string[] environmentVariables, bool wait, string[] args) {
        Directory.SetCurrentDirectory(startDir);
        int val = MsbBindings.ForkExec(path, args, environmentVariables, wait);
        if (val < 0)
            throw new Exception(Errno.GetStringErr(val));
        else
            return val;
    }

    /// <summary>
    /// Mounts /proc. /proc is recommended if you want your OS to like, do stuff.
    /// </summary>
    public static void MountProcFS() {
        Directory.CreateDirectory("/proc", UnixFileMode.UserRead 
                | UnixFileMode.UserWrite 
                | UnixFileMode.UserExecute 
                | UnixFileMode.GroupRead 
                | UnixFileMode.GroupWrite 
                | UnixFileMode.GroupExecute 
                | UnixFileMode.OtherRead 
                | UnixFileMode.OtherWrite 
                | UnixFileMode.OtherExecute); // C# moment

        int val = MsbBindings.Mount("proc", "/proc", "proc", 0, "");
        if (val != 0)
            throw new Exception(Errno.GetStringErr(val));
    }

    /// <summary>
    /// Unmounts /proc.
    /// </summary>
    public static void UnmountProcFS() {
        int val = MsbBindings.Unmount("/proc",0);
        if (val != 0)
            throw new Exception(Errno.GetStringErr(val));
    }

    /// <summary>
    /// Mounts /dev. /dev is also recommended if you want your OS to like, do stuff.
    /// </summary>
    public static void MountDevFS() {
        int val = MsbBindings.Mount("udev", "/dev", "devtmpfs", 2, ""); // 2 = MS_NOSUID
        if (val != 0)
            throw new Exception(Errno.GetStringErr(val));
    }

    /// <summary>
    /// Unmounts /dev.
    /// </summary>
    public static void UnmountDevFS() {
        int val = MsbBindings.Unmount("/dev",0);
        if (val != 0)
            throw new Exception(Errno.GetStringErr(val));
    }
}