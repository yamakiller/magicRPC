using System;

namespace IDL
{
    public class NamespaceGolangCode
    {
        public static string CreateSpaceCode(ParseNamespace namespaceInterface)
        {
            return "package " + namespaceInterface.GetName() + "\n\n";
        }
    }
}