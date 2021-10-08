//magic RPC Client Automatic generation interface:2021-10-07 11:37:42

package login

import (
	"errors"
	rpclient "github.com/yamakiller/magicRPC/client"
	rpcc "github.com/yamakiller/magicRPC/rpc"
)


type testInterface struct {
}

func (t *testInterface) getNameInterface(broker *rpclient.Broker, sign Signin, compressType rpcc.MRPC_PACKAGE_COMPRESS, timeout int) (*UserInfo, error) {
	const funcID = 122943573
	builder := flatbuffers.NewBuilder(512)
	//input flatbuffer code for SigninFB class

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
	req := NewRequest(funcID, compressType, builder.FinishedBytes())
	resp, err := broker.SyncCall(req, rpcc.MRPC_PT_FLATBUFF)
	if err != nil{
		return nil, err
	}

	if resp._responseStatus != rpcc.RS_OK {
		return nil, errors.New("rpc sync call error, function:getName")
	}

	retUserInfoFB := GetRootAsUserInfoFB(resp.Pop(), 0)
	var retUserInfo UserInfo
	//input flatbuffer code for %UserInfoFB class

	retUserInfo.name = string(retUserInfoFB.Name())
	retUserInfo.age = retUserInfoFB.Age()
	for i := 0;i < retUserInfoFB.SexLength();i++ {
		retUserInfo.sex = append(retUserInfo.sex,	retUserInfoFB.Sex(i))
	}

	for i := 0;i < retUserInfoFB.WidgetLength();i++ {
		var widgetTmp Widget
		var widget WidgetFB
		if ok := retUserInfoFB.Widget(&widget,i);!ok {
			return errors.New("deserialize retUserInfoFB=>Widget fail")
		}

		widgetTmp.value1 = widget.Value1()
		widgetTmp.value2 = string(widget.Value2())
		retUserInfo.widget = append(retUserInfo.widget,widgetTmp)
	}

	for i := 0;i < retUserInfoFB.AddressLength();i++ {
		retUserInfo.address = append(retUserInfo.address,	string(retUserInfoFB.Address(i)))
	}


	return &retUserInfo, nil
}
