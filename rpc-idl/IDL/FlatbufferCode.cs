using System;
using System.Collections.Generic;

namespace IDL
{
    public class FlatbufferCode
    {
        static Dictionary<string, string> m_flatbufferVariable = new Dictionary<string, string>();
        public static string CreateFlatbufferCode(ParseStruct structInterface)
        {
            m_flatbufferVariable["bool"] = "bool";
            m_flatbufferVariable["int8"] = "byte";
            m_flatbufferVariable["uint8"] = "ubyte";
            m_flatbufferVariable["int16"] = "short";
            m_flatbufferVariable["uint16"] = "ushort";
            m_flatbufferVariable["int"] = "int";
            m_flatbufferVariable["uint"] = "uint";
            m_flatbufferVariable["int32"] = "long";
            m_flatbufferVariable["uint32"] = "ulong";
            m_flatbufferVariable["float32"] = "float";
            m_flatbufferVariable["float64"] = "double";
            m_flatbufferVariable["string"] = "string";

            string strs = "table " + structInterface.StructName + "FB {\n";

            for (int i = 0; i < structInterface.MemberAttrs.Count; i++)
            {
                MemberAttr v = structInterface.MemberAttrs[i];
                string typeName = getFlatbufferVariable(v.TypeName);
                if (typeName == null)
                {
                    ParseStruct pstruct = (ParseStruct)Vars.GetStruct(v.TypeName);
                    if (pstruct != null)
                    {
                        typeName = pstruct.GetName() + "FB";
                    }
                }

                if (typeName != null)
                {
                    if (v.IsArray)
                    {
                        if (v.IsClass)
                            strs += "\t" + v.VarName + ":[" + v.TypeName + "FB];\n";
                        else
                            strs += "\t" + v.VarName + ":[" + v.TypeName + "];\n";
                    }
                    else
                    {
                        if (v.IsClass)
                            strs += "\t" + v.VarName + ":" + v.TypeName + "FB;\n";
                        else
                            strs += "\t" + v.VarName + ":" + v.TypeName + ";\n";
                    }
                }
                else
                    throw new System.Exception("create flatbuffer file is faild, message name: " +
                    structInterface.StructName + ", type:" + v.TypeName + " is not exits!");
            }
            strs += "}\n\n";

            return strs;
        }

        static string getFlatbufferVariable(string name)
        {
            try
            {
                return m_flatbufferVariable[name];
            }
            catch (System.Collections.Generic.KeyNotFoundException)
            {
                return null;
            }
        }
    }



}