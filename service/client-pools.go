package service

import (
	"reflect"

	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicRPC/rpc"
)

//Pool client object pool
type clientPools struct {
	_core                  *actors.Core
	_delegateBinderAgent   *delegateBinder
	_delegateResponseAgent func(resp *Response, protocol rpc.MRPC_PACKAGE_PROTOCOL) error
}

//Get 返回一个连接对象
func (slf *clientPools) get() *client {
	c := &client{
		Box:    *boxs.SpawnBox(nil),
		_dget:  slf._delegateBinderAgent.get,
		_dresp: slf._delegateResponseAgent,
	}

	slf._core.New(func(pid *actors.PID) actors.Actor {
		c.Box.WithPID(pid)
		c.Box.Register(reflect.TypeOf(&Request{}), c.onRequest)
		return &c.Box
	}, actors.PriorityNomal)

	return c
}

//Put 释放客户端对象
func (slf *clientPools) put(c *client) {
	c.Shutdown()
}
