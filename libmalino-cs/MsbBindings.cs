using System.Runtime.InteropServices;

namespace libmalino;

public class MsbBindings
{
    [DllImport("libmsb", EntryPoint = "msb_sync")]
    public static extern void Sync();

    [DllImport("libmsb", EntryPoint = "msb_reboot")]
    public static extern long Reboot(uint cmd);

    [DllImport("libmsb", EntryPoint = "msb_write")]
    public static extern long Write(uint fd, string buf, ulong count);
}
