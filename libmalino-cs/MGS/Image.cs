using System;
using System.Collections.Generic;
using System.Text;

namespace libmalino.MGS
{
    /// <summary>
    /// Represents a raster image.
    /// </summary>
    public abstract class Image
    {
        /// <summary>
        /// The raw data of the image. This array holds all of the pixel
        /// values of the raster image.
        /// </summary>
        public int[] RawData;

        /// <summary>
        /// The width of the image.
        /// </summary>
        public uint Width { get; protected set; }

        /// <summary>
        /// The height of the image.
        /// </summary>
        public uint Height { get; protected set; }

        /// <summary>
        /// The color depth of each pixel of the image - i.e, the amount
        /// of bits per each pixel.
        /// </summary>
        public ColorDepth Depth { get; protected set; }

        /// <summary>
        /// Initializes a new instance of <see cref="Image"/> class.
        /// </summary>
        /// <param name="width">The width of the image.</param>
        /// <param name="height">The height of the image.</param>
        /// <param name="color">The color depth of each pixel.</param>
        protected Image(uint width, uint height, ColorDepth color)
        {
            Width = width;
            Height = height;
            Depth = color;
        }
        public void resize(uint NewW, uint NewH)
        {
            RawData = ScaleImage(this, (int)NewW, (int)NewH);
        }

        private int[] ScaleImage(Image image, int newWidth, int newHeight)
        {
            int[] rawData = image.RawData;
            int width = (int)image.Width;
            int height = (int)image.Height;
            int[] array = new int[newWidth * newHeight];
            int num = (width << 16) / newWidth + 1;
            int num2 = (height << 16) / newHeight + 1;
            for (int i = 0; i < newHeight; i++)
            {
                for (int j = 0; j < newWidth; j++)
                {
                    int num3 = j * num >> 16;
                    int num4 = i * num2 >> 16;
                    array[i * newWidth + j] = rawData[num4 * width + num3];
                }
            }

            return array;
        }
    }

    /// <summary>
    /// Supported image formats.
    /// </summary>
    public enum ImageFormat
    {
        BMP
    }
}
