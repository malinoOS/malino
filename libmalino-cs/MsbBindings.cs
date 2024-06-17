using System.Runtime.InteropServices;
using System.Text;

namespace libmalino;

public class MsbBindings
{
    [DllImport("libmsb", EntryPoint = "msb_sync")]
    public static extern void Sync();

    [DllImport("libmsb", EntryPoint = "msb_reboot")]
    public static extern long Reboot(uint cmd);

    [DllImport("libmsb", EntryPoint = "msb_write")]
    public static extern long Write(uint fd, string buf);

    [DllImport("libmsb", EntryPoint = "msb_read")]
    public static extern long Read(uint fd, StringBuilder buf, ulong count);
}
