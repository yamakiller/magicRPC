
namespace IDL
{
    public class NamespaceCSharpCode
    {
        public static string CreateSpaceCode(ParseNamespace namespaceInterface)
        {
            return "package " + namespaceInterface.GetName();
        }
    }
}