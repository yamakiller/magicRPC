package client

import "github.com/yamakiller/magicRPC/rpc"

func NewRequest(funcID uint,
	compressType rpc.MRPC_PACKAGE_COMPRESS,
	timeout uint32,
	payload []byte) *Request {
	request := &Request{
		_compressType: compressType,
		_func:         funcID,
		_timeout:      timeout,
	}

	if payload != nil {
		request._payload = make([]byte, len(payload))
		copy(request._payload, payload)
	}

	return request
}

type Request struct {
	_compressType   rpc.MRPC_PACKAGE_COMPRESS
	_responseStatus rpc.RESPONSE_STATUS
	_payload        []byte
	_func           uint
	_nonblock       bool
	_timestamp      uint32
	_timeout        uint32
	_sequeNum       uint32
}

func (r *Request) Push(args []byte) {
	r._payload = args
}

func (r *Request) Pop() []byte {
	return r._payload
}

func (r *Request) bindFunc(fun uint) {
	r._func = fun
}

func (r *Request) GetArgsNum() int {
	return 1
}

func (r *Request) GetCallFuncID() uint {
	return r._func
}

func (r *Request) GetTimestamp() uint32 {
	return r._timestamp
}

func (r *Request) WithSequence(seque uint32) {
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

func (r *Request) WithNonblock(isNonblock bool) {
	r._nonblock = isNonblock
}

func (r *Request) GetCompressType() rpc.MRPC_PACKAGE_COMPRESS {
	return r._compressType
}

func (r *Request) WithCompressType(typ rpc.MRPC_PACKAGE_COMPRESS) {
	r._compressType = typ
}
