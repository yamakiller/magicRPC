
namespace IDL
{
    public class NamespaceCppCode
    {
        public static string CreateSpaceCode(ParseNamespace namespaceInterface)
        {
            return "package " + namespaceInterface.GetName();
        }
    }
}