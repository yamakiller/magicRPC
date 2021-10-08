using System;
using System.Text.RegularExpressions;
using System.Collections.Generic;
using Libs;


using size_t = System.UIntPtr;

namespace IDL
{
    public class ParseInterface : IBParse
    {
        int m_funcIdx = 0;
        string m_interfaceName;
        string m_filePath;
        List<FunctionAttr> m_functions = new List<FunctionAttr>();

        public string InterfaceName { get { return m_interfaceName; } }

        public List<FunctionAttr> Functions { get { return m_functions; } }
        public bool Parse(string filePath, string name, string bodys)
        {
            this.m_interfaceName = name;
            this.m_filePath = filePath;

            string[] memberAttrList = bodys.Split(";");
            if (memberAttrList.Length < 1)
            {
                throw new System.Exception("parse service member attr is failed, " + bodys);
            }

            foreach (string attr in memberAttrList)
            {
                if (attr.Length > 1)
                {
                    FunctionAttr funcAttr = new FunctionAttr(m_filePath + this.m_interfaceName, attr);
                    m_functions.Insert(m_funcIdx++, funcAttr);
                    //m_functionList[] = funcAttr;
                }
            }

            return true;
        }

        public string GetName()
        {
            return m_interfaceName;
        }

        public string CreateCodeForLanguage(ELanguage language)
        {
            return "";
        }

        public string CreateServerInterfaceForLanguage(ELanguage language)
        {
            string codeText;
            switch (language)
            {
                case ELanguage.CL_SHARP: codeText = InterfaceCSharpCode.CreateServerInterfaceCode(this); break;
                case ELanguage.CL_GOLANG: codeText = InterfaceGolangCode.CreateServerInterfaceCode(this); break;
                case ELanguage.CL_CPP: codeText = InterfaceCppCode.CreateServerInterfaceCode(this); break;
                case ELanguage.CL_JAVA: codeText = InterfaceJavaCode.CreateServerInterfaceCode(this); break;
                default:
                    throw new System.Exception("language is not exits!!");
            }
            return codeText;
        }

        public string CreateServerServiceForLanguage(ELanguage language)
        {
            string codeText;
            switch (language)
            {
                case ELanguage.CL_SHARP: codeText = InterfaceCSharpCode.CreateServerServiceCode(this); break;
                case ELanguage.CL_GOLANG: codeText = InterfaceGolangCode.CreateServerServiceCode(this); break;
                case ELanguage.CL_CPP: codeText = InterfaceCppCode.CreateServerServiceCode(this); break;
                case ELanguage.CL_JAVA: codeText = InterfaceJavaCode.CreateServerServiceCode(this); break;
                default:
                    throw new System.Exception("language is not exits!!");
            }
            return codeText;
        }
        public string CreateClientInterfaceForLanguage(ELanguage language)
        {
            string codeText;
            switch (language)
            {
                case ELanguage.CL_SHARP: codeText = InterfaceCSharpCode.CreateClientInterfaceCode(this); break;
                case ELanguage.CL_GOLANG: codeText = InterfaceGolangCode.CreateClientInterfaceCode(this); break;
                case ELanguage.CL_CPP: codeText = InterfaceCppCode.CreateClientInterfaceCode(this); break;
                case ELanguage.CL_JAVA: codeText = InterfaceJavaCode.CreateClientInterfaceCode(this); break;
                default:
                    throw new System.Exception("language is not exits!!");
            }
            return codeText;
        }

        public string CreateClientServiceForLanguage(ELanguage language)
        {
            string codeText;
            switch (language)
            {
                case ELanguage.CL_SHARP: codeText = InterfaceCSharpCode.CreateClientServiceCode(this); break;
                case ELanguage.CL_GOLANG: codeText = InterfaceGolangCode.CreateClientServiceCode(this); break;
                case ELanguage.CL_CPP: codeText = InterfaceCppCode.CreateClientServiceCode(this); break;
                case ELanguage.CL_JAVA: codeText = InterfaceJavaCode.CreateClientServiceCode(this); break;
                default:
                    throw new System.Exception("language is not exits!!");
            }
            return codeText;
        }
    }

    public class FunctionArg
    {
        string m_typeName;
        string m_varName;

        public string TypeName { get { return m_typeName; } }

        public string VarName { get { return m_varName; } }
        public FunctionArg(string type, string var)
        {
            m_typeName = type;
            m_varName = var;

            if (Vars.GetStruct(m_typeName) == null)
            {
                throw new System.Exception("parse type error, is not exist type:" + m_typeName);
            }
        }
    }

    public class FunctionAttr
    {

        string m_funcName;

        size_t m_funcHash;

        FunctionArg m_retValue;

        FunctionArg m_funcArgMap;

        public string FuncName { get { return m_funcName; } }

        public size_t FuncHash { get { return m_funcHash; } }

        public FunctionArg FuncArgMap { get { return m_funcArgMap; } }

        public FunctionArg FuncReturn { get { return m_retValue; } }
        public FunctionAttr(string fileRoutePath, string functlp)
        {
            string fromatFunctlp = Regex.Replace(functlp, @"[\(\)]", " ");
            fromatFunctlp = Regex.Replace(fromatFunctlp, @"\,", "");
            fromatFunctlp = Regex.Replace(fromatFunctlp, @"^\s*", "");
            fromatFunctlp = Regex.Replace(fromatFunctlp, @"\s*$", "");

            string[] funcTlpList = fromatFunctlp.Split(" ");


            if (funcTlpList.Length < 4 && funcTlpList.Length % 2 == 0)
            {
                throw new System.Exception("parse function arguments is failed, " + functlp);
            }
            m_retValue = new FunctionArg(funcTlpList[0], "ret" + funcTlpList[0]);

            m_funcName = funcTlpList[1];
            m_funcHash = (size_t)HashEncoder.Hash(System.Text.Encoding.UTF8.GetBytes(fileRoutePath + m_funcName));

            for (int i = 2; i < funcTlpList.Length; i += 2)
            {
                m_funcArgMap = new FunctionArg(funcTlpList[i], funcTlpList[i + 1]);
            }
        }
    }
}