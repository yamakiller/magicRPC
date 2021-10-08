package service

import (
	"fmt"
	"sync/atomic"

	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicRPC/rpc"
	"github.com/yamakiller/magicRPC/rpc/rpctime"
)

//sender 发送者
//type sender func(sock interface{}, data interface{}) error

//type closer func(sock interface{})

type client struct {
	boxs.Box
	_dget   delegateGet
	_dresp  func(resp *Response, protocol rpc.MRPC_PACKAGE_PROTOCOL) error
	_socket int32 //套接字
	_ref    int32 //引用计数器
}

func (c *client) onRequest(context *boxs.Context) {
	var (
		err       error
		cb        *delegateHandle
		timespace uint32
	)
	request := context.Message().(*Request)
	util.AssertEmpty(request, fmt.Sprintf("on request event deserialization error:%+v", context.Message()))

	if request._timeout > 0 {
		timespace = rpctime.Click()
		if request._timestamp+request._timeout >= timespace {
			request.WithStatus(rpc.RS_TIMEOUT)
			goto fail_lable
		}
	}

	cb = c._dget(request.GetCallFuncID())
	if cb == nil {
		context.Error("call fail:function %d undefine", request.GetCallFuncID())
		request.WithStatus(rpc.RS_FAILD)
		goto fail_lable
	}

	if err = cb._func(context, request); err != nil {
		context.Error("call %s return error:%v", cb._name, err)
		request.WithStatus(rpc.RS_FAILD)
		goto fail_lable
	}

	if request._timeout > 0 {
		timespace = rpctime.Click()
		if request._timestamp+request._timeout >= timespace {
			context.Error("call %s timeout", cb._name)
		}
	}

	return
fail_lable:

	if err = c._dresp(NewReponse(request, nil), rpc.MRPC_PT_FLATBUFF); err != nil {
		context.Error("response timeout error:%v", err)
	}
}

//Ref  获取引用计数器
func (c *client) ref() int32 {
	return atomic.LoadInt32(&c._ref)
}

//RefInc 引用计数器+1
func (c *client) refInc() int32 {
	return atomic.AddInt32(&c._ref, 1)
}

//RefDec 引用计数器-1
func (c *client) refDec() int32 {
	return atomic.AddInt32(&c._ref, -1)
}
