using System;
using System.Drawing;
using System.IO;
using System.IO.MemoryMappedFiles;
using System.Runtime.CompilerServices;
using System.Text;
namespace libmalino.MGS;
    public class Canvas
    {
        public Mode Mode { get; set; }
        public byte[] buffer { get; set; }
        MemoryMappedFile nmap;
        public Canvas()
        {
            var b = File.ReadAllBytes("/dev/fb0");
            nmap = MemoryMappedFile.CreateFromFile("/dev/fb0", FileMode.Open, null, b.Length);

            buffer = new byte[b.Length];
            var res = ((ScreenResolution)(b.Length / 4)).ToString().Replace("_","");
            var WH = res.Split('x');
            Mode = new Mode(uint.Parse(WH[0]),uint.Parse(WH[1]),ColorDepth.ColorDepth32);
        }
        #region displayClasses
        public void Display()
        {
            var stream = nmap.CreateViewStream();
            stream.Write(buffer,0,buffer.Length);
            stream.Seek(1, SeekOrigin.Current);
            stream.Close();
        }
        public void Disable()
        {
            malinoIO.ClearScreen();
        }
        #endregion
        
        #region DrawingClasses
        public void DrawPoint(Color col, int x, int y)
        {
            int index = (x + (y * (int)Mode.Width)) * 4;
            buffer[index] = col.B;
            buffer[index + 1] = col.G;
            buffer[index + 2] = col.R;
            buffer[index + 3] = 255;
        }
        public void DrawFilledRectangle(Color col,int x,int y,int w,int h)
        {
            for (int i = x; i < x+w; i++)
            {
                for (int j = y; j < y+h; j++)
                {
                    DrawPoint(col,i,j);
                }
            }
        }
        public void Clear(Color col)
        {
            DrawFilledRectangle(col,0,0,(int)Mode.Width,(int)Mode.Height);
        }
        public Color GetPointColor(int x, int y)
        {
            int index = (x + (y * (int)Mode.Width)) * 4;
            return Color.FromArgb(buffer[index + 3],buffer[index + 2],buffer[index + 1],buffer[index]);
        }
        public void DrawLine(Color color, int x1, int y1, int x2, int y2)
        {
            int dx = Math.Abs(x2 - x1);
            int dy = Math.Abs(y2 - y1);
            int sx = x1 < x2 ? 1 : -1;
            int sy = y1 < y2 ? 1 : -1;
            int err = dx - dy;

            while (true)
            {

                DrawPoint(color,x1, y1);

                if (x1 == x2 && y1 == y2)
                    break;

                int e2 = 2 * err;
                if (e2 > -dy)
                {
                    err -= dy;
                    x1 += sx;
                }
                if (e2 < dx)
                {
                    err += dx;
                    y1 += sy;
                }
            }
        }
        static string ASC16Base64 = "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AAB4zMzM2MzGxsbMAAAAAAAA/sbGwMDAwMDAwAAAAAAAAAAA/mxsbGxsbGwAAAAAAAAA/sZgMBgwYMb+AAAAAAAAAAAAftjY2NjYcAAAAAAAAAAAZmZmZmZ8YGDAAAAAAAAAAHbcGBgYGBgYAAAAAAAAAH4YPGZmZjwYfgAAAAAAAAA4bMbG/sbGbDgAAAAAAAA4bMbGxmxsbGzuAAAAAAAAHjAYDD5mZmZmPAAAAAAAAAAAAH7b29t+AAAAAAAAAAAAAwZ+29vzfmDAAAAAAAAAHDBgYHxgYGAwHAAAAAAAAAB8xsbGxsbGxsYAAAAAAAAAAP4AAP4AAP4AAAAAAAAAAAAYGH4YGAAA/wAAAAAAAAAwGAwGDBgwAH4AAAAAAAAADBgwYDAYDAB+AAAAAAAADhsbGBgYGBgYGBgYGBgYGBgYGBgYGNjY2HAAAAAAAAAAABgYAH4AGBgAAAAAAAAAAAAAdtwAdtwAAAAAAAAAOGxsOAAAAAAAAAAAAAAAAAAAAAAAABgYAAAAAAAAAAAAAAAAAAAAGAAAAAAAAAAADwwMDAwM7GxsPBwAAAAAANhsbGxsbAAAAAAAAAAAAABw2DBgyPgAAAAAAAAAAAAAAAAAfHx8fHx8fAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==";
        static MemoryStream ASC16FontMS = new MemoryStream(Convert.FromBase64String(ASC16Base64));
        public void DrawString(string s, int x, int y, Color col)
        {

            string[] lines = s.Split('\n');
            for (int l = 0; l < lines.Length; l++)
            {
                for (int c = 0; c < lines[l].Length; c++)
                {
                    int offset = (Encoding.ASCII.GetBytes(lines[l][c].ToString())[0] & 0xFF) * 16;
                    ASC16FontMS.Seek(offset, SeekOrigin.Begin);
                    byte[] fontbuf = new byte[16];
                    ASC16FontMS.Read(fontbuf, 0, fontbuf.Length);

                    for (int i = 0; i < 16; i++)
                    {
                        for (int j = 0; j < 8; j++)
                        {
                            if ((fontbuf[i] & (0x80 >> j)) != 0)
                            {
                                if (!(x + c * 8 > Mode.Width))
                                {
                                    DrawPoint(col,(int)((x + j) + (c * 8)), (int)(y + i + (l * 16)));
                                }
                            }
                        }
                    }
                }
            }

        }
        public void DrawRectangle(int left, int top, int width, int height, Color color)
        {

            DrawLine(color,left,top,width,top);
            DrawLine(color,width,top,width,height);
            DrawLine(color,left,height,width,height);
            DrawLine(color,left,top,left,height);
        }
        public void DrawCircle(Color color, int xCenter, int yCenter, int radius)
        {
            int x = radius;
            int y = 0;
            int e = 0;

            while (x >= y)
            {
                DrawPoint(color, xCenter + x, yCenter + y);
                DrawPoint(color, xCenter + y, yCenter + x);
                DrawPoint(color, xCenter - y, yCenter + x);
                DrawPoint(color, xCenter - x, yCenter + y);
                DrawPoint(color, xCenter - x, yCenter - y);
                DrawPoint(color, xCenter - y, yCenter - x);
                DrawPoint(color, xCenter + y, yCenter - x);
                DrawPoint(color, xCenter + x, yCenter - y);

                y++;
                if (e <= 0)
                {
                    e += (2 * y) + 1;
                }
                if (e > 0)
                {
                    x--;
                    e -= (2 * x) + 1;
                }
            }
        }
        public void DrawFilledCircle(Color color, int x0, int y0, int radius)
        {
            int x = radius;
            int y = 0;
            int xChange = 1 - (radius << 1);
            int yChange = 0;
            int radiusError = 0;

            while (x >= y)
            {
                for (int i = x0 - x; i <= x0 + x; i++)
                {

                    DrawPoint(color, i, y0 + y);
                    DrawPoint(color, i, y0 - y);
                }
                for (int i = x0 - y; i <= x0 + y; i++)
                {
                    DrawPoint(color, i, y0 + x);
                    DrawPoint(color, i, y0 - x);
                }

                y++;
                radiusError += yChange;
                yChange += 2;
                if ((radiusError << 1) + xChange > 0)
                {
                    x--;
                    radiusError += xChange;
                    xChange += 2;
                }
            }
        }
        public void DrawImageAlpha(Bitmap image, int x, int y)
        {
            for (int i = 0; i < image.Width; i++)
            {
                for (int j = 0; j < image.Height; j++)
                {
                    Color color = Color.FromArgb(image.RawData[i + j * image.Width]);
                    if (color.A > 0)
                    {
                        if(color.A < 255)
                        {
                            color = AlphaBlend(color, GetPointColor(x+i,y+j), color.A);
                        }
                        DrawPoint(color,x+i,y+j);

                    }
                        
                }
            }
        }
        public void DrawEllipse(Color color, int xCenter, int yCenter, int xR, int yR)
        {
            int a = 2 * xR;
            int b = 2 * yR;
            int b1 = b & 1;
            int dx = 4 * (1 - a) * b * b;
            int dy = 4 * (b1 + 1) * a * a;
            int err = dx + dy + (b1 * a * a);
            int e2;
            int y = 0;
            int x = xR;
            a *= 8 * a;
            b1 = 8 * b * b;

            while (x >= 0)
            {
                DrawPoint(color, xCenter + x, yCenter + y);
                DrawPoint(color, xCenter - x, yCenter + y);
                DrawPoint(color, xCenter - x, yCenter - y);
                DrawPoint(color, xCenter + x, yCenter - y);
                e2 = 2 * err;
                if (e2 <= dy) { y++; err += dy += a; }
                if (e2 >= dx || 2 * err > dy) { x--; err += dx += b1; }
            }
        }
        public void DrawFilledEllipse(Color color, int xCenter, int yCenter, int yR, int xR)
        {
            for (int y = -yR; y <= yR; y++)
            {
                for (int x = -xR; x <= xR; x++)
                {
                    if ((x * x * yR * yR) + (y * y * xR * xR) <= yR * yR * xR * xR)
                    {
                        DrawPoint(color, xCenter + x, yCenter + y);
                    }
                }
            }
        }
        public void DrawArc(int x, int y, int width, int height, Color color, int startAngle = 0, int endAngle = 360)
        {
            if (width == 0 || height == 0)
            {
                return;
            }

            for (double angle = startAngle; angle < endAngle; angle += 0.5)
            {
                double angleRadians = Math.PI * angle / 180;
                int IX = (int)(width * Math.Cos(angleRadians));
                int IY = (int)(height * Math.Sin(angleRadians));
                DrawPoint(color, x + IX, y + IY);
            }
        }
        public void DrawPolygon(Color color, params Point[] points)
        {
            // Using an array of points here is better than using something like a Dictionary of ints.
            if (points.Length < 3)
            {
                throw new ArgumentException("A polygon requires more than 3 points.");
            }

            for (int i = 0; i < points.Length - 1; i++)
            {
                var pointA = points[i];
                var pointB = points[i + 1];
                DrawLine(color, pointA.X, pointA.Y, pointB.X, pointB.Y);
            }

            var firstPoint = points[0];
            var lastPoint = points[^1];
            DrawLine(color, firstPoint.X, firstPoint.Y, lastPoint.X, lastPoint.Y);
        }
        public void DrawSquare(Color color, int x, int y, int size)
        {
            DrawRectangle(color, x, y, size, size);
        }
        public void DrawRectangle(Color color, int x, int y, int width, int height)
        {
            int xa = x;
            int ya = y;

            int xb = x + width;
            int yb = y;

            int xc = x;
            int yc = y + height;

            int xd = x + width;
            int yd = y + height;

            DrawLine(color, xa, ya, xb, yb);

            DrawLine(color, xa, ya, xc, yc);

            DrawLine(color, xb, yb, xd, yd);

            DrawLine(color, xc, yc, xd, yd);
        }
        public void DrawTriangle(Color color, int v1x, int v1y, int v2x, int v2y, int v3x, int v3y)
        {
            DrawLine(color, v1x, v1y, v2x, v2y);
            DrawLine(color, v1x, v1y, v3x, v3y);
            DrawLine(color, v2x, v2y, v3x, v3y);
        }
        public void DrawImage(Image image, int x, int y, bool preventOffBoundPixels = true)
        {
            Color color;
            if (preventOffBoundPixels)
            {
                var maxWidth = Math.Min(image.Width, (int)Mode.Width - x);
                var maxHeight = Math.Min(image.Height, (int)Mode.Height - y);
                for (int xi = 0; xi < maxWidth; xi++)
                {
                    for (int yi = 0; yi < maxHeight; yi++)
                    {
                        color = Color.FromArgb(image.RawData[xi + (yi * image.Width)]);
                        DrawPoint(color, x + xi, y + yi);
                    }
                }
            }
            else
            {
                for (int xi = 0; xi < image.Width; xi++)
                {
                    for (int yi = 0; yi < image.Height; yi++)
                    {
                        color = Color.FromArgb(image.RawData[xi + (yi * image.Width)]);
                        DrawPoint(color, x + xi, y + yi);
                    }
                }
            }
        }
        public void DrawArray(Color[] colors, int x, int y, int width, int height)
        {
            for (int X = 0; X < width; X++)
            {
                for (int Y = 0; Y < height; Y++)
                {
                    DrawPoint(colors[Y * width + X], x + X, y + Y);
                }
            }
        }
        public static Color AlphaBlend(Color to, Color from, byte alpha)
        {
            byte R = (byte)(((to.R * alpha) + (from.R * (255 - alpha))) >> 8);
            byte G = (byte)(((to.G * alpha) + (from.G * (255 - alpha))) >> 8);
            byte B = (byte)(((to.B * alpha) + (from.B * (255 - alpha))) >> 8);
            return Color.FromArgb(R, G, B);
        }
        #endregion
    }