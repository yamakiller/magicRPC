using System;
using System.IO;
using Libs;
using IDL;

namespace rpc_idl
{
    class Program
    {
        static void Main(string[] args)
        {
            Libs.Args pargs = new Args(args);
            Console.WriteLine("111111");
            string inFile = pargs.Get("-in", null);
            if (inFile == null)
            {
                Console.WriteLine("请输入需要解析的文件,参数-in=xxxxx");
                return;
            }
            Console.WriteLine("2222222");
            string outFilePath = pargs.Get("-out", null);
            if (outFilePath == null)
            {
                Console.WriteLine("请输入输出文件目录,参数-out=....");
                return;
            }
            Console.WriteLine("3333333");
            string language = pargs.Get("-language", "go");
            IDL.Builder pbuilder = new Builder(language);
            pbuilder.ParseFile = inFile;
            pbuilder.OutputFilePath = outFilePath;
            pbuilder.FileName = Path.GetFileNameWithoutExtension(inFile);
            Console.WriteLine("444444444");
            pbuilder.StartParse();
            Console.WriteLine("sdfsfsfsdf");
        }
    }
}
