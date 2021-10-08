using System;
using System.Collections.Generic;

namespace Libs
{
    public class Args
    {
        Dictionary<string, string> m_args = new Dictionary<string, string>();
        public Args(string[] args)
        {
            foreach (string arg in args)
            {
                string[] m = arg.Split(':');
                if (m.Length <= 1)
                {
                    m_args.Add(arg, arg);
                    m_args[arg] = arg;
                    continue;
                }

                m_args[m[0]] = m[1];
            }
        }

        public string Get(string key, string def)
        {
            try
            {
                return m_args[key];
            }
            catch (System.Collections.Generic.KeyNotFoundException)
            {
                return def;
            }

        }
    }

}