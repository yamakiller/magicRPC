using System;
using System.Diagnostics;

namespace Libs
{
    public class Proc
    {
        /// <summary>
        /// 调用Exe核心代码
        /// </summary>
        /// <param name="exeFileName"></param>
        /// <param name="args"></param>
        public static void RunExe(string exeFileName, string args = "")
        {
            try
            {

                Process p = new Process();

                p.StartInfo = new ProcessStartInfo(exeFileName, args);

                p.StartInfo.Arguments = args;
                //p.StartInfo.WorkingDirectory = @"C:\MinGW\msys\1.0\home\61125\wfdb-10.6.1\build\bin";
                p.StartInfo.UseShellExecute = false;

                p.StartInfo.RedirectStandardOutput = true;

                //p.StartInfo.RedirectStandardInput = true;

                p.StartInfo.RedirectStandardError = true;

                p.StartInfo.CreateNoWindow = false;
                //绑定事件
                p.OutputDataReceived += new DataReceivedEventHandler(p_OutputDataReceived);
                p.ErrorDataReceived += p_ErrorDataReceived;

                p.Start();
                p.BeginOutputReadLine();//开始读取输出数据
                p.BeginErrorReadLine();//开始读取错误数据，重要！
                p.WaitForExit();
                p.Close();
            }
            catch (Exception e)
            {
                Console.WriteLine(e);
                throw;
            }
        }

        private static void p_ErrorDataReceived(object sender, DataReceivedEventArgs e)
        {
            var output = e.Data;
            if (output != null)
            {

                Console.WriteLine(output);
            }
        }

        private static void p_OutputDataReceived(object sender, DataReceivedEventArgs e)
        {
            Console.WriteLine(e.Data);
        }
    }
}