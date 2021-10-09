using System;
using System.Collections.Generic;
using IDL;
using Libs;

namespace IDL
{
    public class ParseStruct : IBParse
    {
        string m_structName;
        List<MemberAttr> m_memberAttrInfo = new List<MemberAttr>();

        public string StructName { get { return m_structName; } }

        public List<MemberAttr> MemberAttrs { get { return m_memberAttrInfo; } }
        public bool Parse(string filename, string name, string bodys)
        {
            if (((int)name[0] < 0x41 && (int)name[0] > 0x51))
            {
                throw new System.Exception("parse meesgae name is failed, The first character of the name is capitalized, message name:" + name);
            }

            m_structName = name;

            string[] memberAttrList = bodys.Split(';');
            if (memberAttrList.Length < 1)
            {
                throw new System.Exception("parse meesgae member attr is failed, " + bodys);
            }

            foreach (string attr in memberAttrList)
            {
                if (string.IsNullOrEmpty(attr))
                    continue;

                string[] member = attr.Split(':');
                if (member.Length <= 1)
                    throw new System.Exception("parse message member flag is failed, " + member[0]);
                int index = int.Parse(member[1]);
                string[] memberFlags = member[0].Split(' ');
                if (memberFlags.Length < 2)
                    throw new System.Exception("parse message member flag is failed, " + member[0]);
                try
                {
                    m_memberAttrInfo.Insert(index, new MemberAttr(memberFlags[0], memberFlags[1]));
                }
                catch (System.ArgumentOutOfRangeException)
                {
                    throw new System.Exception("struct index error:" + index);
                }
            }
            return true;
        }

        public string GetName()
        {
            return m_structName;
        }

        public string CreateCodeForLanguage(ELanguage language)
        {
            string codeText = "";

            switch (language)
            {
                case ELanguage.CL_SHARP: codeText = StructCSharppCode.CreateServerCode(this); break;
                case ELanguage.CL_GOLANG: codeText = StructGolangCode.CreateServerCode(this); break;
                case ELanguage.CL_CPP: codeText = StructCppCode.CreateServerCode(this); break;
                case ELanguage.CL_JAVA: codeText = StructJavaCode.CreateServerCode(this); break;
                default:
                    throw new System.Exception("language is not exits!!");
            }

            return codeText;
        }


        public string CreateServerInterfaceForLanguage(ELanguage language)
        {
            return "";
        }

        public string CreateServerServiceForLanguage(ELanguage language)
        {
            return "";
        }

        public string CreateClientInterfaceForLanguage(ELanguage language)
        {
            return "";
        }

        public string CreateClientServiceForLanguage(ELanguage language)
        {
            return "";
        }



        static int deserializeRecursive = -1;
        public static string CreateDeserializeCodeForFlatbuffer(ParseStruct structInfo, string varName, string fbName)
        {
            string strs = "";
            string spacesStr = "";
            List<MemberAttr> memberAttrInfo = structInfo.MemberAttrs;
            string structName = structInfo.StructName;
            deserializeRecursive++;
            for (int j = 0; j <= deserializeRecursive; j++)
                spacesStr += "\t";

            for (int i = 0; i < memberAttrInfo.Count; i++)
            {
                string iterName = memberAttrInfo[i].VarName;
                string iterType = memberAttrInfo[i].TypeName;

                //for (int j = 0; j <= deserializeRecursive; j++)
                //     speces += "\t";

                if (memberAttrInfo[i].IsArray)
                {
                    if (memberAttrInfo[i].IsClass)
                    {
                        string tmpName = iterType.ToLower() + "Tmp";
                        strs += spacesStr + "for i := 0;i < " + fbName + "." + StringTo.ToUpper(iterName) + "Length();i++ {\n";
                        strs += spacesStr + "\tvar " + tmpName + " " + iterType + "\n";
                        strs += spacesStr + "\tvar " + iterName + " " + iterType + "FB\n";

                        strs += spacesStr + "\tif ok := " + fbName + "." + StringTo.ToUpper(iterName) + "(&" + iterName + ",i);!ok {\n";
                        strs += spacesStr + "\t\treturn nil, errors.New(\"deserialize " + fbName + "=>" + StringTo.ToUpper(iterName) + " fail\")\n";
                        strs += spacesStr + "\t}\n\n";
                        strs += ParseStruct.CreateDeserializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(iterType), tmpName, iterName);
                        strs += spacesStr + "\t" + varName + "." + iterName + " = append(" + varName + "." + iterName + "," + tmpName + ")\n";
                    }
                    else
                    {
                        strs += spacesStr + "for i := 0;i < " + fbName + "." + StringTo.ToUpper(iterName) + "Length();i++ {\n";
                        strs += spacesStr + "\t" + varName + "." + iterName + " = append(" + varName + "." + iterName + ",";
                        if (memberAttrInfo[i].IsString)
                            strs += spacesStr + "string(" + fbName + "." + StringTo.ToUpper(iterName) + "(i)))\n";
                        else
                            strs += spacesStr + fbName + "." + StringTo.ToUpper(iterName) + "(i))\n";
                    }
                    strs += spacesStr + "}\n\n";
                }
                else if (memberAttrInfo[i].IsClass)
                {
                    strs += ParseStruct.CreateDeserializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(iterType), varName + "." + iterName,
                    fbName + "." + StringTo.ToUpper(iterName) + "()");
                }
                else
                {
                    if (memberAttrInfo[i].IsString)
                        strs += spacesStr + varName + "." + iterName + " = string(" + fbName + "." + StringTo.ToUpper(iterName) + "())\n";
                    else
                        strs += spacesStr + varName + "." + iterName + " = " + fbName + "." + StringTo.ToUpper(iterName) + "()\n";
                }
            }

            deserializeRecursive--;
            return strs;
        }

        static int serializeRecursive = -1;
        public static string CreateSerializeCodeForFlatbuffer(ParseStruct structInfo, string varName)
        {
            string classStr = "";
            string stringStr = "";
            string argsStr = "";
            string spacesStr = "";
            List<MemberAttr> memberAttrInfo = structInfo.MemberAttrs;
            string structName = structInfo.StructName;

            serializeRecursive++;
            for (int j = 0; j <= serializeRecursive; j++)
                spacesStr += "\t";



            for (int i = 0; i < memberAttrInfo.Count; i++)
            {
                string iterName = memberAttrInfo[i].VarName;
                string iterType = memberAttrInfo[i].TypeName;

                if (memberAttrInfo[i].IsArray)
                {
                    if (memberAttrInfo[i].IsClass)
                    {
                        classStr += "\n";
                        classStr += spacesStr + "//class " + iterType + " vector serialize\n";
                        classStr += spacesStr + iterName + "PosArray := make([]flatbuffers.UOffsetT, len(" + varName + "." + iterName + "))\n";
                        classStr += spacesStr + "for i," + iterName + " := range " + varName + "." + iterName + " {\n";
                        classStr += spacesStr + ParseStruct.CreateSerializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(iterType), iterName);
                        classStr += spacesStr + "\t" + iterName + "PosArray[i] = " + iterName + "Pos\n";

                    }
                    else if (memberAttrInfo[i].IsString)
                    {
                        classStr += "\n";
                        classStr += spacesStr + "//string " + iterName + " vector serialize\n";
                        classStr += spacesStr + iterName + "PosArray := make([]flatbuffers.UOffsetT, len(" + varName + "." + iterName + "))\n";
                        classStr += spacesStr + "for i," + iterName + " := range " + varName + "." + iterName + " {\n";
                        classStr += spacesStr + "\t" + iterName + "PosArray[i] = builder.CreateString(" + iterName + ")\n";
                    }
                    else
                    {
                        string endOffset = varName + StringTo.ToUpper(iterName) + "EndOffset";
                        classStr += "\n";
                        classStr += spacesStr + "//" + iterType + " " + iterName + " vector serialize\n";
                        classStr += spacesStr + structName + "FBStart" + StringTo.ToUpper(iterName) + "Vector(builder, len(" + varName + "." + iterName + "))\n";
                        classStr += spacesStr + "for _," + iterName + " := range " + varName + "." + iterName + " {\n";
                        classStr += spacesStr + "\tbuilder.Prepend" + StringTo.ToUpper(iterType) + "(" + iterName + ")\n";
                        classStr += spacesStr + "}\n";
                        classStr += spacesStr + endOffset + " := builder.EndVector(len(" + varName + "." + iterName + "))\n\n";
                        argsStr += spacesStr + structName + "FBAdd" + StringTo.ToUpper(iterName) + "(builder," + endOffset + ")\n";
                    }

                    if (memberAttrInfo[i].IsClass || memberAttrInfo[i].IsString)
                    {
                        classStr += spacesStr + "}\n";

                        string endOffset = varName + StringTo.ToUpper(iterName) + "EndOffset";

                        classStr += spacesStr + structName + "FBStart" + StringTo.ToUpper(iterName) + "Vector(builder, len(" + iterName + "PosArray))\n";
                        classStr += spacesStr + "for _," + iterName + "Offset := range " + iterName + "PosArray {\n";
                        classStr += "\t\tbuilder.PrependUOffsetT(" + iterName + "Offset)\n";
                        classStr += "\t}\n";
                        classStr += spacesStr + endOffset + " := builder.EndVector(len(" + iterName + "PosArray))\n\n";
                        argsStr += spacesStr + structName + "FBAdd" + StringTo.ToUpper(iterName) + "(builder," +
                        varName + StringTo.ToUpper(iterName) + "EndOffset)\n";
                    }



                }
                else if (memberAttrInfo[i].IsClass)
                {
                    classStr += "\n";
                    classStr += "//class " + iterType + " " + iterName + " serialize\n";
                    classStr += ParseStruct.CreateSerializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(iterType), iterName);
                    argsStr += spacesStr + structName + "FBAdd" + StringTo.ToUpper(iterName) + "(builder," + iterName + "Pos" + ")\n";
                }
                else if (memberAttrInfo[i].IsString)
                {
                    string stringPos = varName + StringTo.ToUpper(iterName) + "StringPos";
                    stringStr += "\n";
                    stringStr += spacesStr + stringPos + " := builder.CreateString(" + varName + "." + iterName + ")\n";
                    argsStr += spacesStr + structName + "FBAdd" + StringTo.ToUpper(iterName) + "(builder," + stringPos + ")\n";
                }
                else
                {
                    argsStr += spacesStr + structName + "FBAdd" + StringTo.ToUpper(iterName) + "(builder," + varName + "." + iterName + ")\n";
                }
            }


            classStr += spacesStr + "//class " + structName + " serialize start\n";
            classStr += stringStr;
            classStr += spacesStr + structName + "FBStart(builder)\n";
            classStr += argsStr;
            classStr += spacesStr + varName + "Pos := " + structName + "FBEnd(builder)\n";
            classStr += spacesStr + "//class " + structName + " serialize end\n\n";
            serializeRecursive--;
            return classStr;
        }
    }
    //增加创建与FlatBuffer的序列化及反序列化代码
}


public class MemberAttr
{
    string m_typeName;
    string m_varName;
    bool m_isArray;
    bool m_isString;
    bool m_isClass;

    public string TypeName { get { return m_typeName; } }

    public string VarName { get { return m_varName; } }
    public bool IsArray { get { return m_isArray; } }
    public bool IsString { get { return m_isString; } }
    public bool IsClass { get { return m_isClass; } }
    public MemberAttr(string type, string member)
    {
        {
            string[] types = type.Split(new char[2] { '[', ']' });
            if (types.Length > 1)
            {
                m_isArray = true;
                m_typeName = types[0];
            }
            else
            {
                m_isArray = false;
                m_typeName = type;
            }
        }

        if (m_typeName == "string")
            m_isString = true;
        else
            m_isString = false;

        if (Vars.GetStruct(m_typeName) != null)
        {
            m_isClass = true;
        }
        else
        {
            if (Vars.GetVariable(m_typeName) == null)
                throw new System.Exception("Idl Incorrect type, type name:" + m_typeName);
            m_typeName = Vars.GetVariable(m_typeName);
            m_isClass = false;
        }
        m_varName = member;
    }


}
