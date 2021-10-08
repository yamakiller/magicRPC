package service

import (
	"github.com/yamakiller/magicRPC/rpc"
)

type Request struct {
	_sock            interface{}
	_networkProtocol rpc.MRPC_NETWORK_PROTOCOL
	_compressType    rpc.MRPC_PACKAGE_COMPRESS
	_responseStatus  rpc.RESPONSE_STATUS
	_playload        []byte
	_func            uint
	_nonblock        bool
	_timestamp       uint32
	_timeout         uint32
	_sequeNum        uint32
}

func (r *Request) bindSock(sock interface{}) {
	r._sock = sock
}

func (r *Request) Sock() interface{} {
	return r._sock
}

func (r *Request) Push(args []byte) {
	r._playload = args
}

func (r *Request) Pop() []byte {
	return r._playload
}

func (r *Request) bindFunc(fun uint) {
	r._func = fun
}

func (r *Request) withNetworkProtocol(netproto rpc.MRPC_NETWORK_PROTOCOL) {
	r._networkProtocol = netproto
}

func (r *Request) GetNet() rpc.MRPC_NETWORK_PROTOCOL {
	return r._networkProtocol
}

func (r *Request) GetArgsNum() int {
	return 1
}

func (r *Request) GetCallFuncID() uint {
	return r._func
}

/*(func (r *Request) GetFunArgList() []byte {
	return r._funcArgs
}*/

func (r *Request) GetTimestamp() uint32 {
	return r._timestamp
}

func (r *Request) withSequence(seque uint32) {
	r._sequeNum = seque
}

func (r *Request) GetSequence() uint32 {
	return r._sequeNum
}

func (r *Request) GetTimeout() uint32 {
	return r._timeout
}

func (r *Request) WithStatus(status rpc.RESPONSE_STATUS) {
	r._responseStatus = status
}

func (r *Request) GetStatus() rpc.RESPONSE_STATUS {
	return r._responseStatus
}

func (r *Request) GetNonblock() bool {
	return r._nonblock
}

func (r *Request) withNonblock(isNonblock bool) {
	r._nonblock = isNonblock
}

func (r *Request) getCompressType() rpc.MRPC_PACKAGE_COMPRESS {
	return r._compressType
}

func (r *Request) withCompressType(typ rpc.MRPC_PACKAGE_COMPRESS) {
	r._compressType = typ
}
