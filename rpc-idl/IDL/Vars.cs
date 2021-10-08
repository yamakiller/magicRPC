using System;
using System.Collections.Generic;

namespace IDL
{
    public class Vars
    {
        static IBParse m_namespace = null;
        static Dictionary<string, IBParse> m_structs = new Dictionary<string, IBParse>();
        static Dictionary<string, IBParse> m_interfaces = new Dictionary<string, IBParse>();
        static Dictionary<string, string> m_variables = new Dictionary<string, string>();
        public static IBParse GetNamespace()
        {
            return m_namespace;
        }

        public static void RegisterNamespace(IBParse space)
        {
            m_namespace = space;
        }


        public static IBParse GetStruct(string name)
        {
            try
            {
                return m_structs[name];
            }
            catch (System.Collections.Generic.KeyNotFoundException)
            {
                return null;
            }
        }

        public static void RegisterStruct(string name, IBParse parseStruct)
        {
            m_structs[name] = parseStruct;
        }

        public static IBParse GetInterface(string name)
        {
            try
            {
                return m_interfaces[name];
            }
            catch (System.Collections.Generic.KeyNotFoundException)
            {
                return null;
            }
        }

        public static void RegisterInterface(string name, IBParse parseInterface)
        {
            m_interfaces[name] = parseInterface;
        }

        public static string GetVariable(string name)
        {
            try
            {
                return m_variables[name];
            }
            catch (System.Collections.Generic.KeyNotFoundException)
            {
                return null;
            }
        }

        public static void RegisterVariable(string idlVarName, string objVarName)
        {
            m_variables.Add(idlVarName, objVarName);
        }

        public static Dictionary<string, IBParse> GetStructs()
        {
            return m_structs;
        }

        public static Dictionary<string, IBParse> GetInterfaces()
        {
            return m_interfaces;
        }
    }
}