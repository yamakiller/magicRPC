namespace IDL
{
    public class ParseNamespace : IBParse
    {
        string m_spacename;
        public bool Parse(string filename, string name, string bodys)
        {
            m_spacename = name;
            return true;
        }

        public string GetName()
        {
            return m_spacename;
        }

        public string CreateCodeForLanguage(ELanguage language)
        {
            string codeText = "";

            switch (language)
            {
                case ELanguage.CL_SHARP: codeText = NamespaceCSharpCode.CreateSpaceCode(this); break;
                case ELanguage.CL_GOLANG: codeText = NamespaceGolangCode.CreateSpaceCode(this); break;
                case ELanguage.CL_CPP: codeText = NamespaceCppCode.CreateSpaceCode(this); break;
                case ELanguage.CL_JAVA: codeText = NamespaceJavaCode.CreateSpaceCode(this); break;
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
    }
}

