using System;
using System.IO;


namespace Libs
{
    public class FileSave
    {
        public static bool Save(string filePath, string data)
        {
            FileStream fs = null;
            if (File.Exists(filePath))
                //
                fs = new FileStream(filePath, FileMode.Truncate);
            else
                fs = new FileStream(filePath, FileMode.OpenOrCreate);
            StreamWriter sw = new StreamWriter(fs);
            sw.Write(data);
            sw.Close();
            fs.Close();
            fs.Dispose();
            return true;
        }
    }
}