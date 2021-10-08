using System;

namespace IDL
{
    public class StructGolangCode
    {
        public static string CreateServerCode(ParseStruct parse)
        {
            string strs = "type " + parse.StructName + " struct {\n";
            foreach (MemberAttr attr in parse.MemberAttrs)
            {
                if (attr.IsArray)
                    strs += "\t" + attr.VarName + "[] " + attr.TypeName + "\n";
                else
                    strs += "\t" + attr.VarName + " " + attr.TypeName + "\n";
            }
            strs += "}\n\n\n";

            return strs;
        }


    }

}