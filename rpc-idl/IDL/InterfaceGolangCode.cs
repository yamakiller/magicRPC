using System;
using System.Collections.Generic;
using Libs;

namespace IDL
{
    public class InterfaceGolangCode
    {
        public static string CreateServerInterfaceCode(ParseInterface parse)
        {

            string interfaceClassName = parse.InterfaceName + "Interface";
            string interfaceClassVar = parse.InterfaceName[0] +
            Md5Helper.Md5(interfaceClassName).ToLower().Replace("-", "");
            interfaceClassVar = interfaceClassVar.ToLower();

            string strs = "var (\n\t" + interfaceClassVar + " *" + interfaceClassName + "\n)\n\n";

            strs += "func Fasten" + parse.InterfaceName + "(psrv *rpcs.Server) {\n";
            strs += "\t" + interfaceClassVar + " = &" + interfaceClassName + "{_implsrv:psrv}\n";
            foreach (FunctionAttr funAttr in parse.Functions)
            {
                strs += "\t" + interfaceClassVar + "._implsrv.BindDelegate(" + funAttr.FuncHash + "," +
                interfaceClassVar + "." + funAttr.FuncName + "Interface," +
                "\"" + parse.InterfaceName + "Service." + funAttr.FuncName + "\")\n";
            }
            strs += "}\n\n";


            strs += "type " + parse.InterfaceName + "Interface struct {\n";
            strs += "\t_implsrv *rpcs.Server\n";
            strs += "}\n\n";

            foreach (FunctionAttr funAttr in parse.Functions)
            {
                strs += functionAttrCode.CreateServerInterfaceCode(funAttr, parse.InterfaceName);
            }

            return strs;
        }

        public static string CreateServerServiceCode(ParseInterface parse)
        {
            string serviceName = parse.InterfaceName + "Service";
            string strs = "type " + serviceName + " struct {\n";
            strs += "\t_implcontext *rpcb.Context\n";
            strs += "}\n\n";

            foreach (FunctionAttr funAttr in parse.Functions)
            {
                strs += functionAttrCode.CreateServerServiceCode(funAttr, parse.InterfaceName);
            }

            return strs;
        }

        public static string CreateClientInterfaceCode(ParseInterface parse)
        {
            string interfaceClassName = StringTo.ToLower(parse.InterfaceName, 0) + "Interface";
            string strs = "type " + interfaceClassName + " struct {\n";
            strs += "}\n\n";
            foreach (FunctionAttr funAttr in parse.Functions)
            {
                strs += functionAttrCode.CreateClientInterfaceCode(funAttr, parse.InterfaceName);
            }

            return strs;
        }

        public static string CreateClientServiceCode(ParseInterface parse)
        {
            string serviceName = parse.InterfaceName + "Service";
            string strs = "type " + serviceName + " struct {\n";
            strs += "\t" + "_brokerImpl *rpclient.Broker\n";
            strs += "\t" + "_compressType rpcc.MRPC_PACKAGE_COMPRESS\n";
            strs += "}\n\n";
            strs += "func (" + serviceName[0].ToString().ToLower() + " *" + serviceName + ") WithBroker(broker *rpclient.Broker) {\n";
            strs += "\t" + serviceName[0].ToString().ToLower() + "._brokerImpl = broker\n";
            strs += "}\n";

            strs += "func (" + serviceName[0].ToString().ToLower() + " *" + serviceName + ") WithCommpressType(compressType rpcc.MRPC_PACKAGE_COMPRESS) {\n";
            strs += "\t" + serviceName[0].ToString().ToLower() + "._compressType = compressType\n";
            strs += "}\n";


            foreach (FunctionAttr funAttr in parse.Functions)
            {
                strs += functionAttrCode.CreateClientServiceCode(funAttr, parse.InterfaceName);
            }

            return strs;
        }
    }

    class functionAttrCode
    {
        public static string CreateServerInterfaceCode(FunctionAttr functionAttrInterface, string inerfaceName)
        {
            string strs = "func (slf *" + inerfaceName + "Interface)  " + functionAttrInterface.FuncName + "Interface(context *rpcb.Context, req *rpcs.Request) error{\n";
            strs += "\t" + functionAttrInterface.FuncArgMap.VarName + "FB := GetRootAs" + functionAttrInterface.FuncArgMap.TypeName + "FB(req.Pop(), 0)\n";
            strs += "\t" + functionAttrInterface.FuncArgMap.VarName + " := " + functionAttrInterface.FuncArgMap.TypeName + "{}\n";

            strs += "\t//input flatbuffer code for " + functionAttrInterface.FuncArgMap.VarName + "FB class\n";

            strs += "" + ParseStruct.CreateDeserializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(functionAttrInterface.FuncArgMap.TypeName),
                        functionAttrInterface.FuncArgMap.VarName, functionAttrInterface.FuncArgMap.VarName + "FB") + "\n\n";

            strs += "\tdeleageService := " + inerfaceName + "Service{_implcontext:context}\n";
            strs += "\t" + functionAttrInterface.FuncReturn.VarName + " := deleageService." + functionAttrInterface.FuncName + "(" + functionAttrInterface.FuncArgMap.VarName + ")\n";


            strs += "\tbuilder := slf._implsrv.MallocFlatBuilder()\n";
            strs += "\tdefer slf._implsrv.FreeFlatBuilder(builder)\n";
            strs += "\t//input flatbuffer code for " + functionAttrInterface.FuncReturn.TypeName + "\n";
            strs += "\t" + ParseStruct.CreateSerializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(functionAttrInterface.FuncReturn.TypeName),
                functionAttrInterface.FuncReturn.VarName);
            strs += "\tbuilder.Finish(" + functionAttrInterface.FuncReturn.VarName + "Pos)\n";

            strs += "\tresp := rpcs.NewReponse(req, builder.FinishedBytes())\n";
            strs += "\tif err := slf._implsrv.Response(resp, rpcc.MRPC_PT_FLATBUFF); err != nil {\n";
            strs += "\t\tcontext.Error(\"response fail:%v\", err)\n";
            strs += "\t}\n";


            strs += "\treturn nil\n}\n\n";

            return strs;
        }

        public static string CreateServerServiceCode(FunctionAttr functionAttrInterface, string inerfaceName)
        {
            string strs = "";
            string funcArgsStrs = "";
            string returnStrs = StringTo.ToLower(functionAttrInterface.FuncReturn.TypeName) + "Ret";
            string serviceName = inerfaceName + "Service";

            FunctionArg v = functionAttrInterface.FuncArgMap;
            funcArgsStrs += v.VarName + " " + v.TypeName;

            strs += "func (" + StringTo.ToLower(functionAttrInterface.FuncReturn.TypeName)[0] + " *" + serviceName + ") " +
            functionAttrInterface.FuncName + "(" + funcArgsStrs + ") " + functionAttrInterface.FuncReturn.TypeName + " {\n";
            strs += "\tvar " + returnStrs + " " + functionAttrInterface.FuncReturn.TypeName + "\n";

            strs += "\t//TODO: 编写函数代码\n\n";
            strs += "\treturn " + returnStrs + "\n";
            strs += "}\n";

            return strs;
        }

        public static string CreateClientInterfaceCode(FunctionAttr functionAttrInterface, string inerfaceName)
        {
            string strs = "";
            string funcArgsStr = "";
            string funcArgsStructStr = "";
            string funcInterfaceName = StringTo.ToLower(inerfaceName, 0) + "Interface";


            FunctionArg v = functionAttrInterface.FuncArgMap;
            funcArgsStr = v.VarName + " " + v.TypeName;
            funcArgsStructStr = v.VarName;


            strs += "func (" + inerfaceName[0].ToString().ToLower() + " *" + funcInterfaceName + ") " +
                       StringTo.ToLower(functionAttrInterface.FuncName, 0) + "Interface(broker" + " *rpclient.Broker, " +
                       funcArgsStr + ", compressType rpcc.MRPC_PACKAGE_COMPRESS, timeout int) (*" + functionAttrInterface.FuncReturn.TypeName + ", error) {\n";
            strs += "\tconst funcID = " + functionAttrInterface.FuncHash + "\n";
            strs += "\tbuilder := flatbuffers.NewBuilder(512)\n";
            strs += "\t//input flatbuffer code for " + v.TypeName + "FB class\n";
            strs += ParseStruct.CreateSerializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(functionAttrInterface.FuncReturn.TypeName),
                                                                         functionAttrInterface.FuncReturn.VarName);
            strs += "\tbuilder.Finish(" + functionAttrInterface.FuncReturn.VarName + "Pos)\n";
            strs += "\treq := NewRequest(funcID, compressType, builder.FinishedBytes())\n";
            strs += "\tresp, err := broker.SyncCall(req, rpcc.MRPC_PT_FLATBUFF)\n";
            strs += "\tif err != nil{\n";
            strs += "\t\treturn nil, err\n";
            strs += "\t}\n\n";
            strs += "\tif resp._responseStatus != rpcc.RS_OK {\n";
            strs += "\t\treturn nil, errors.New(\"rpc sync call error, function:" + functionAttrInterface.FuncName + "\")\n";
            strs += "\t}\n\n";

            strs += "\t" + functionAttrInterface.FuncReturn.VarName + "FB := GetRootAs" +
                               functionAttrInterface.FuncReturn.TypeName + "FB(resp.Pop(), 0)\n";
            strs += "\tvar " + functionAttrInterface.FuncReturn.VarName + " " + functionAttrInterface.FuncReturn.TypeName + "\n";
            strs += "\t//input flatbuffer code for %" + functionAttrInterface.FuncReturn.TypeName + "FB class\n\n";
            strs += ParseStruct.CreateDeserializeCodeForFlatbuffer((ParseStruct)Vars.GetStruct(functionAttrInterface.FuncReturn.TypeName),
                             functionAttrInterface.FuncReturn.VarName, functionAttrInterface.FuncReturn.VarName + "FB") + "\n";
            strs += "\treturn &" + functionAttrInterface.FuncReturn.VarName + ", nil\n";
            strs += "}\n";

            return strs;
        }

        public static string CreateClientServiceCode(FunctionAttr functionAttrInterface, string inerfaceName)
        {
            string serviceName = StringTo.ToUpper(inerfaceName, 0) + "Service";
            string interfaceName = StringTo.ToLower(inerfaceName, 0) + "Interface";

            string selfName = inerfaceName[0].ToString().ToLower();
            FunctionArg v = functionAttrInterface.FuncArgMap;
            string funcArgsStr = v.VarName + " " + v.TypeName;

            string strs = "func (" + selfName + " *" + serviceName + ") " + StringTo.ToUpper(functionAttrInterface.FuncName, 0) + "(";
            strs += funcArgsStr + ", timeout int" + ") (*" + functionAttrInterface.FuncReturn.TypeName + ", error) {\n";
            strs += "\treturn " + interfaceName + "{}." + StringTo.ToLower(functionAttrInterface.FuncName, 0) + "Interface(" + selfName + "._brokerImpl" +
                    ", " + v.VarName + ", " + selfName + "._compressType, timeout)\n";
            strs += "}\n\n";
            return strs;
        }
    }
}