using System;
using System.Collections.Generic;
using System.IO;

namespace libmalino;

/// <summary>
/// Core functions.
/// </summary>
#pragma warning disable IDE1006
#pragma warning disable CS8981
public class malino {
#pragma warning restore CS8981
#pragma warning restore IDE1006
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
    /// Spawns a process on the system.
    ///
    /// Returns: if wait=true, the return code, if wait=false, 0. This will throw an exception of the errno message if an error happened while spawning the process.
    /// </summary>
    public static int SpawnProcess(string path, string startDir, string[] environmentVariables, bool wait, string[] args) {
        Directory.SetCurrentDirectory(startDir);
        List<string> argTemp = [path];
        argTemp.AddRange(args);
        
        #pragma warning disable CS8625
        argTemp.Add(null);
        List<string> envp = [.. environmentVariables, null];
        #pragma warning restore CS8625

        int val = MsbBindings.ForkExec(path, [.. argTemp], [.. envp], wait);
        if (val < 0)
            throw new Exception(Errno.GetStringErr(-val));
        else
            return val;
    }

    /// <summary>
    /// Mounts /proc. /proc is recommended if you want your OS to do stuff.
    /// </summary>
    public static void MountProcFS() {
        // C#: hmm yes, let's make a warning about creating a folder, and not about the linux system call in this function
        // I know this isn't supported by windows, stfu dotnet
        #pragma warning disable CA1416
        Directory.CreateDirectory("/proc", UnixFileMode.UserRead 
                | UnixFileMode.UserWrite 
                | UnixFileMode.UserExecute 
                | UnixFileMode.GroupRead 
                | UnixFileMode.GroupWrite 
                | UnixFileMode.GroupExecute 
                | UnixFileMode.OtherRead 
                | UnixFileMode.OtherWrite 
                | UnixFileMode.OtherExecute); // C# moment
        #pragma warning restore CA1416

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
    /// Mounts /dev. /dev is also recommended if you want your OS to do stuff.
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