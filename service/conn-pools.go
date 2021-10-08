package service

import (
	"sync"

	"github.com/yamakiller/magicNet/netboxs"
)

//Pool connection 对象池
type connPools struct {
	_bfsize int
	_wqsize int
	_pps    *sync.Pool
	_once   sync.Once
}

//Get 返回一个连接对象
func (slf *connPools) Get() netboxs.Connect {
	slf._once.Do(func() {
		slf._pps = &sync.Pool{
			New: func() interface{} {
				c := &Connection{netboxs.BTCPConn{ReadBufferSize: slf._bfsize,
					WriteBufferSize: slf._bfsize,
					WriteQueueSize:  slf._wqsize,
				}}
				return c
			},
		}
	})

	return slf._pps.Get().(netboxs.Connect)
}

//Put 释放客户端对象
func (slf *connPools) Put(c netboxs.Connect) {
	slf._pps.Put(c)
}
