using System;
using System.IO;
using System.Text.RegularExpressions;
using System.Collections.Generic;
using Libs;

namespace IDL
{
    public class Builder
    {

        string m_inFilePath = ".";
        string m_outFilePath = ".";
        string m_fileName;
        string m_moduleFilePath;
        ELanguage m_language;

        public string ParseFile { set { m_inFilePath = value; } }
        public string OutputFilePath { set { m_outFilePath = value; } }
        public string FileName { set { m_fileName = value; } }

        public Builder(string language)
        {
            switch (language.ToLower().Trim())
            {
                case "go": m_language = ELanguage.CL_GOLANG; golangVariableInit(); break;
                case "cpp": m_language = ELanguage.CL_CPP; cppVariableInit(); break;
                case "csharp": m_language = ELanguage.CL_SHARP; csharpVariableInit(); break;
                case "java": m_language = ELanguage.CL_JAVA; javaVariableInit(); break;
                default: m_language = ELanguage.CL_GOLANG; golangVariableInit(); break;
            }
        }

        public void StartParse()
        {
            string text = "";
            try
            {
                text = System.IO.File.ReadAllText(m_inFilePath);
            }
            catch (System.Exception e)
            {
                Console.WriteLine("read file fail:" + e.Message);
                return;
            }
            this.parse(text);
        }

        bool parse(string data)
        {
            data = Regex.Replace(data, @"\/\/[^\n]*", "");
            data = Regex.Replace(data, @"[\n\r\t]", "");
            data = Regex.Replace(data, @"\s{2,}", "");

            string[] classes = Regex.Split(data, @"[@\}]");

            if (classes.Length == 0)
            {
                throw new System.Exception("parse classes is failed, no class struct!!");
            }

            //

            foreach (string c in classes)
            {
                string[] symbolFlag = c.Split('{');
                if (symbolFlag.Length != 2)
                {
                    continue;
                }

                string[] symbolAttr = symbolFlag[0].Split(":");
                if (symbolAttr.Length != 2)
                {
                    throw new Exception("parse symbol  attr is failed,  symbol missing :, " + symbolFlag[0]);
                }

                IBParse idlParse;
                switch (symbolAttr[0])
                {
                    case Symbol.Struct:
                        idlParse = new ParseStruct();
                        break;
                    case Symbol.Interface:
                        idlParse = new ParseInterface();
                        break;
                    case Symbol.Namespace:
                        idlParse = new ParseNamespace();
                        break;
                    default:
                        throw new Exception("parse symbol attr is error,  symbol: " + symbolAttr[0]);
                }

                if (idlParse.Parse(m_fileName, symbolAttr[1].Trim(), symbolFlag[1].Trim()))
                {
                    switch (symbolAttr[0])
                    {
                        case Symbol.Struct:
                            Vars.RegisterStruct(idlParse.GetName(), idlParse);
                            break;
                        case Symbol.Interface:
                            Vars.RegisterInterface(idlParse.GetName(), idlParse);
                            break;
                        case Symbol.Namespace:
                            Vars.RegisterNamespace(idlParse);
                            break;
                    }
                }
            }

            //TODDO:解释代码修
            createCode();
            return true;
        }

        void createCode()
        {
            if (Vars.GetNamespace() == null)
            {
                throw new System.Exception("undefine namespace");
            }

            string serverCodeInterface, serverCodeService;
            string structCode;
            string flatbufferCode;

            serverCodeService = "//magic RPC Automatic generation service:" + DateTime.Now.ToString("yyyy-MM-dd hh:mm:ss") + "\n\n";
            serverCodeService += Vars.GetNamespace().CreateCodeForLanguage(m_language);

            serverCodeInterface = "//magic RPC Automatic generation interface:" + DateTime.Now.ToString("yyyy-MM-dd hh:mm:ss") + "\n\n";
            serverCodeInterface += Vars.GetNamespace().CreateCodeForLanguage(m_language);

            serverCodeInterface += "import (\n";
            serverCodeInterface += "\trpcs \"github.com/yamakiller/magicRPC/service\"\n";
            serverCodeInterface += "\trpcc \"github.com/yamakiller/magicRPC/rpc\"\n";
            serverCodeInterface += "\trpcb \"github.com/yamakiller/magicLibs/boxs\"\n";
            serverCodeInterface += "\tflatbuffers \"github.com/google/flatbuffers/go\"\n";
            serverCodeInterface += ")\n\n\n";


            serverCodeService += "import (\n";
            serverCodeService += "\trpcb \"github.com/yamakiller/magicLibs/boxs\"\n";
            serverCodeService += ")\n\n\n";

            string clientCodeInterface, clientCodeService;

            clientCodeInterface = "//magic RPC Client Automatic generation interface:" + DateTime.Now.ToString("yyyy-MM-dd hh:mm:ss") + "\n\n";
            clientCodeInterface += Vars.GetNamespace().CreateCodeForLanguage(m_language);
            clientCodeInterface += "import (\n";
            clientCodeInterface += "\t\"errors\"\n";
            clientCodeInterface += "\trpclient \"github.com/yamakiller/magicRPC/client\"\n";
            clientCodeInterface += "\trpcc \"github.com/yamakiller/magicRPC/rpc\"\n";
            clientCodeInterface += ")\n\n\n";

            clientCodeService = "//magic RPC Client Automatic generation service:" + DateTime.Now.ToString("yyyy-MM-dd hh:mm:ss") + "\n\n";
            clientCodeService += Vars.GetNamespace().CreateCodeForLanguage(m_language);
            clientCodeService += "import (\n";
            clientCodeService += "\trpclient \"github.com/yamakiller/magicRPC/client\"\n";
            clientCodeService += "\trpcc \"github.com/yamakiller/magicRPC/rpc\"\n";
            clientCodeService += ")\n\n\n";

            /*foreach(k, v; idlInerfaceList)
            {
                serverCodeInterface ~= v.createServerCodeForInterface(CODE_LANGUAGE.CL_DLANG);
                serverCodeService ~= v.createServerCodeForService(CODE_LANGUAGE.CL_DLANG);

                clientCodeInterface ~= v.createClientCodeForInterface(CODE_LANGUAGE.CL_DLANG);
                clientCodeService ~= v.createClientCodeForService(CODE_LANGUAGE.CL_DLANG);
            }*/
            //Interface Code make
            foreach (KeyValuePair<string, IBParse> pair in Vars.GetInterfaces())
            {
                serverCodeInterface += pair.Value.CreateServerInterfaceForLanguage(m_language);
                serverCodeService += pair.Value.CreateServerServiceForLanguage(m_language);

                clientCodeInterface += pair.Value.CreateClientInterfaceForLanguage(m_language);
                clientCodeService += pair.Value.CreateClientServiceForLanguage(m_language);

            }

            //Struct Code make
            structCode = "//magic RPC Automatic generation message:" + DateTime.Now.ToString("yyyy-MM-dd hh:mm:ss") + "\n\n";
            structCode += Vars.GetNamespace().CreateCodeForLanguage(m_language);
            foreach (KeyValuePair<string, IBParse> pair in Vars.GetStructs())
            {
                structCode += pair.Value.CreateCodeForLanguage(m_language);
            }

            //flatbufferCode = "//magic RPC Automatic generation fbs:" + DateTime.Now.ToString("yyyy-MM-dd hh:mm:ss") + "\n\n";
            flatbufferCode = "namespace " + Vars.GetNamespace().GetName() + ";\n\n";
            flatbufferCode += "attribute \"priority\";\n\n";
            foreach (KeyValuePair<string, IBParse> pair in Vars.GetStructs())
            {
                flatbufferCode += FlatbufferCode.CreateFlatbufferCode((ParseStruct)pair.Value);
            }

            if (!Directory.Exists(m_outFilePath))
                Directory.CreateDirectory(m_outFilePath);

            if (!Directory.Exists(m_outFilePath + "/Server"))
                Directory.CreateDirectory(m_outFilePath + "/Server");

            if (!Directory.Exists(m_outFilePath + "/Server/" + Vars.GetNamespace().GetName()))
                Directory.CreateDirectory(m_outFilePath + "/Server/" + Vars.GetNamespace().GetName());

            if (!Directory.Exists(m_outFilePath + "/Client"))
                Directory.CreateDirectory(m_outFilePath + "/Client");

            if (!Directory.Exists(m_outFilePath + "/Client/" + Vars.GetNamespace().GetName()))
                Directory.CreateDirectory(m_outFilePath + "/Client/" + Vars.GetNamespace().GetName());


            if (!Directory.Exists(m_outFilePath + "/Message"))
                Directory.CreateDirectory(m_outFilePath + "/Message");

            if (!Directory.Exists(m_outFilePath + "/Message/" + Vars.GetNamespace().GetName()))
                Directory.CreateDirectory(m_outFilePath + "/Message/" + Vars.GetNamespace().GetName());


            FileSave.Save(m_outFilePath + "/Server/" + Vars.GetNamespace().GetName() + "/" + m_fileName + "Interface" + extension(), serverCodeInterface);
            FileSave.Save(m_outFilePath + "/Server/" + Vars.GetNamespace().GetName() + "/" + m_fileName + "Service" + extension(), serverCodeService);
            FileSave.Save(m_outFilePath + "/Message/" + Vars.GetNamespace().GetName() + "/" + m_fileName + "Message" + extension(), structCode);
            FileSave.Save(m_outFilePath + "/Client/" + Vars.GetNamespace().GetName() + "/" + m_fileName + "Interface" + extension(), clientCodeInterface);
            FileSave.Save(m_outFilePath + "/Client/" + Vars.GetNamespace().GetName() + "/" + m_fileName + "Service" + extension(), clientCodeService);

            FileSave.Save(m_outFilePath + "/" + m_fileName + ".fbs", flatbufferCode);

            Proc.RunExe("flatc.exe", "--go -o  " + m_outFilePath + "/Message/ " + m_outFilePath + "/" + m_fileName + ".fbs");
        }

        void golangVariableInit()
        {
            Vars.RegisterVariable("bool", "bool");
            Vars.RegisterVariable("byte", "byte");
            Vars.RegisterVariable("ubyte", "uint8");
            Vars.RegisterVariable("short", "int16");
            Vars.RegisterVariable("ushort", "uint16");
            Vars.RegisterVariable("int", "int32");
            Vars.RegisterVariable("uint", "uint32");
            Vars.RegisterVariable("int32", "int32");
            Vars.RegisterVariable("uint32", "uint32");
            Vars.RegisterVariable("long", "int64");
            Vars.RegisterVariable("ulong", "uint64");
            Vars.RegisterVariable("float", "float32");
            Vars.RegisterVariable("double", "float64");
            Vars.RegisterVariable("string", "string");
        }

        void cppVariableInit()
        {
            //Vars.RegisterVariable("bool", "bool");
            //Vars.RegisterVariable("byte", "int8_t");
        }

        void csharpVariableInit()
        {

        }

        void javaVariableInit()
        {

        }


        string extension()
        {
            switch (m_language)
            {
                case ELanguage.CL_GOLANG:
                    return ".go";
                case ELanguage.CL_SHARP:
                    return ".cs";
                case ELanguage.CL_CPP:
                    return ".cpp";
                case ELanguage.CL_JAVA:
                    return ".java";
                default:
                    return ".go";
            }
        }
    }
}

