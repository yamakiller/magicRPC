using System.Text;

namespace Libs
{
    public class StringTo
    {
        public static string ToUpper(string src, int pos = 0)
        {
            StringBuilder sbcap = new StringBuilder(src);

            sbcap[pos] = sbcap[pos].ToString().ToUpper()[pos];
            return sbcap.ToString();
        }

        public static string ToLower(string src, int pos = 0)
        {
            StringBuilder sbcap = new StringBuilder(src);

            sbcap[pos] = sbcap[pos].ToString().ToLower()[pos];
            return sbcap.ToString();
        }
    }
}