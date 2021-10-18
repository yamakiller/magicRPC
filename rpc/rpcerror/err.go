package rpcerror

import "errors"

var (
	ErrProtocol             = errors.New("protocol error")
	ErrProtocolVersion      = errors.New("protocol version error")
	ErrProtocolHeaderValid  = errors.New("protocol header valid")
	ErrProtocolDataOverflow = errors.New("protocol data overflow")
	ErrProtocolKeepalive    = errors.New("protocol keepalive error")
	ErrIsProtocolKeepalive  = errors.New("protocol keepalive ok")
	ErrUndefineCompressMode = errors.New("undefine compress mode")
)

var (
	ErrNameNeedConfig       = errors.New("need config name informat")
	ErrLoggerNeedConfig     = errors.New("need config logger")
	ErrBrokerBusy           = errors.New("broker busy")
	ErrBrokerRequestTimeOut = errors.New("broker request timeout")
)
