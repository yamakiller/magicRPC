//magic RPC Automatic generation interface:2021-10-07 11:37:42

package login

import (
	rpcs "github.com/yamakiller/magicRPC/service"
	rpcc "github.com/yamakiller/magicRPC/rpc"
	rpcb "github.com/yamakiller/magicLibs/boxs"
	flatbuffers "github.com/google/flatbuffers/go"
)


var (
	t2b0d349901cf2849abf73afb3c151183 *TestInterface
)

func FastenTest(psrv *rpcs.Server) {
	t2b0d349901cf2849abf73afb3c151183 = &TestInterface{_implsrv:psrv}
	t2b0d349901cf2849abf73afb3c151183._implsrv.BindDelegate(122943573,t2b0d349901cf2849abf73afb3c151183.getNameInterface,"TestService.getName")
}

type TestInterface struct {
	_implsrv *rpcs.Server
}

func (slf *TestInterface)  getNameInterface(context *rpcb.Context, req *rpcs.Request) error{
	signFB := GetRootAsSigninFB(req.Pop(), 0)
	sign := Signin{}
	//input flatbuffer code for signFB class
	sign.name = string(signFB.Name())
	sign.pwd = string(signFB.Pwd())


	deleageService := TestService{_implcontext:context}
	retUserInfo := deleageService.getName(sign)
	builder := slf._implsrv.MallocFlatBuilder()
	defer slf._implsrv.FreeFlatBuilder(builder)
	//input flatbuffer code for UserInfo
	
	//int16 sex vector serialize
	UserInfoFBStartSexVector(builder, len(retUserInfo.sex))
	for _,sex := range retUserInfo.sex {
		builder.PrependInt16(sex)
	}
	retUserInfoSexEndOffset := builder.EndVector(len(retUserInfo.sex))


	//class Widget vector serialize
	widgetPosArray := make([]flatbuffers.UOffsetT, len(retUserInfo.widget))
	for i,widget := range retUserInfo.widget {
			//class Widget serialize start

		widgetValue2StringPos := builder.CreateString(widget.value2)
		WidgetFBStart(builder)
		WidgetFBAddValue1(builder,widget.value1)
		WidgetFBAddValue2(builder,widgetValue2StringPos)
		widgetPos := WidgetFBEnd(builder)
		//class Widget serialize end

		widgetPosArray[i] = widgetPos
	}
	UserInfoFBStartWidgetVector(builder, len(widgetPosArray))
	for _,widgetOffset := range widgetPosArray {
		builder.PrependUOffsetT(widgetOffset)
	}
	retUserInfoWidgetEndOffset := builder.EndVector(len(widgetPosArray))


	//string address vector serialize
	addressPosArray := make([]flatbuffers.UOffsetT, len(retUserInfo.address))
	for i,address := range retUserInfo.address {
		addressPosArray[i] = builder.CreateString(address)
	}
	UserInfoFBStartAddressVector(builder, len(addressPosArray))
	for _,addressOffset := range addressPosArray {
		builder.PrependUOffsetT(addressOffset)
	}
	retUserInfoAddressEndOffset := builder.EndVector(len(addressPosArray))

	//class UserInfo serialize start

	retUserInfoNameStringPos := builder.CreateString(retUserInfo.name)
	UserInfoFBStart(builder)
	UserInfoFBAddName(builder,retUserInfoNameStringPos)
	UserInfoFBAddAge(builder,retUserInfo.age)
	UserInfoFBAddSex(builder,retUserInfoSexEndOffset)
	UserInfoFBAddWidget(builder,retUserInfoWidgetEndOffset)
	UserInfoFBAddAddress(builder,retUserInfoAddressEndOffset)
	retUserInfoPos := UserInfoFBEnd(builder)
	//class UserInfo serialize end

	builder.Finish(retUserInfoPos)
	resp := rpcs.NewReponse(req, builder.FinishedBytes())
	if err := slf._implsrv.Response(resp, rpcc.MRPC_PT_FLATBUFF); err != nil {
		context.Error("response fail:%v", err)
	}
	return nil
}

