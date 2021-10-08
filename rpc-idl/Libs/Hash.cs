using System;

namespace Libs
{

    public class HashEncoder
    {
        public static long Hash(byte[] digest)
        {
            long h = 0;
            for (int i = 0; i < digest.Length; i++)
            {
                h = (h << 4) + digest[i];
                long g = h & 0xF0000000L;
                if (g != 0)
                    h ^= g >> 24;
                h &= ~g;
            }
            return h & 0x7FFFFFFF;
        }
    }

}