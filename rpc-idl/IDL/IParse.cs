namespace IDL
{
    public interface IBParse
    {
        bool Parse(string filename, string name, string bodys);
        string GetName();

        string CreateCodeForLanguage(ELanguage language);

        string CreateServerInterfaceForLanguage(ELanguage language);

        string CreateServerServiceForLanguage(ELanguage language);

        string CreateClientInterfaceForLanguage(ELanguage language);

        string CreateClientServiceForLanguage(ELanguage language);
    }
}