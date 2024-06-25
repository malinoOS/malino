using System;
using System.IO;
using System.IO.MemoryMappedFiles;
using System.Threading;
using System.Runtime.InteropServices;
using System.Drawing;


namespace libmalino.FirstMalino;
    [StructLayout(LayoutKind.Sequential, Pack = 1)]
    struct FrameBufferInfo 
    {
        public int fd;
        public IntPtr data;
        public uint w;
        public uint h;
        public uint bpp;
        public uint line_length;
    }

    class FrameBuffer
    {
        [StructLayout(LayoutKind.Sequential)]
        public struct fb_fix_screeninfo {
            [MarshalAs(UnmanagedType.ByValArray, SizeConst = 16)] public byte[] id;
            [MarshalAs(UnmanagedType.U4)] public uint smem_start;
            [MarshalAs(UnmanagedType.U4)] public uint smem_len;
            [MarshalAs(UnmanagedType.U4)] public uint type;
            [MarshalAs(UnmanagedType.ByValArray, SizeConst = 36)] public byte[] stuff;
        };

        [StructLayout(LayoutKind.Sequential)]
        public struct fb_var_screeninfo {
            public int xres;
            public int yres;
            public int xres_virtual;
            public int yres_virtual;
            public int xoffset;
            public int yoffset;
            public int bits_per_pixel;
            [MarshalAs(UnmanagedType.ByValArray, SizeConst = 132)] public byte[] stuff;
        };


        [DllImport("libc", EntryPoint = "close", SetLastError = true)]
        public static extern int Close(int handle);

        [DllImport("libc", EntryPoint = "ioctl", SetLastError = true)]
        public static extern int Ioctl(int handle, uint request, ref fb_fix_screeninfo capability);


        [DllImport("libc", EntryPoint = "ioctl", SetLastError = true)]
        public static extern int Ioctl(int handle, uint request, ref fb_var_screeninfo capability);

        [DllImport("libc", EntryPoint = "open", SetLastError = true)]
        public static extern int Open(string path, uint flag);

        [DllImport("libc", EntryPoint = "mmap", SetLastError = true)]
        public static extern int Mmap(
            [MarshalAs(UnmanagedType.U4)] uint addr,
            [MarshalAs(UnmanagedType.U4)] uint length,
            [MarshalAs(UnmanagedType.I4)] int prot,
            [MarshalAs(UnmanagedType.I4)] int flags,
            [MarshalAs(UnmanagedType.I4)] int fdes,
            [MarshalAs(UnmanagedType.I4)] int offset
        );

        [DllImport("libc", EntryPoint = "munmap", SetLastError = true)]
        public static extern int Munmap(
            [MarshalAs(UnmanagedType.I4)] int addr,
            [MarshalAs(UnmanagedType.U4)] uint length
        );

        public static T ByteArrayToStructure<T>(byte[] bytes) where T: struct 
        {
            T stuff;
            GCHandle handle = GCHandle.Alloc(bytes, GCHandleType.Pinned);
            try
            {
                stuff = (T)Marshal.PtrToStructure(handle.AddrOfPinnedObject(), typeof(T))!;
            }
            finally
            {
                handle.Free();
            }
            return stuff;
        }
    }