
namespace IDL
{
    public class NamespaceJavaCode
    {
        public static string CreateSpaceCode(ParseNamespace namespaceInterface)
        {
            return "package " + namespaceInterface.GetName() + ";";
        }
    }
}