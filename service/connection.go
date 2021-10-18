package service

import (
	"time"

	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/netboxs"
	"github.com/yamakiller/magicRPC/rpc"
	"github.com/yamakiller/magicRPC/rpc/rpctime"
)

type Connection struct {
	netboxs.BTCPConn
}

//Keepalive 心跳时间
func (c *Connection) Keepalive() time.Duration {
	return time.Duration(_keepaliveTimeSecond) * time.Second
}

//Ping 心跳检测
func (c *Connection) Ping() {
	//发送心跳请求
	resp := &Response{
		_sock:            c.Socket(),
		_networkProtocol: rpc.MRPC_NETWORK_PROTOCOL_TCP,
		_responseStatus:  rpc.RS_OK,
		_func:            0,
		_compressType:    rpc.MRPC_CT_NO,
		_timestamp:       0,
		_timeout:         0,
		_nonblock:        false,
		_sequeNum:        0,
	}

	c.Push(resp)
}

//UnSeria 反序列化
func (c *Connection) UnSeria() (interface{}, error) {
	pk, err := rpc.Decoding(c.Reader())
	if err != nil {
		return nil, err
	}

	//kleepalive data
	if pk.Func() == 0 {
		return nil, nil
	}

	return pk, nil
}

//Seria 序列化
func (c *Connection) Seria(msg interface{}) error {
	var err error
	resp := msg.(*Response)
	util.AssertEmpty(resp, "seria fail, send to data not Response struct")
	timespace := rpctime.Click()
	if resp.GetStatus() != rpc.RS_OK && resp._timeout > 0 {
		if resp._timestamp+resp._timeout > timespace {
			resp.WithStatus(rpc.RS_TIMEOUT)
		}
	}
	h := rpc.Header{}
	h.Init(rpc.MRPC_PT_FLATBUFF, resp._sequeNum, resp._compressType, resp._nonblock, resp._func)
	packageSize := h.GetPackageSize()
	if len(resp._playload) > (rpc.MRPC_PACKAGE_MAX*2 - c.Writer().Buffered() - packageSize) {
		if c.Writer().Buffered() > 0 {
			if err = c.Writer().Flush(); err != nil {
				return err
			}
		}
	}
	if _, err = rpc.Encoding(c.Writer(),
		rpc.MRPC_PT_FLATBUFF,
		resp._sequeNum,
		resp._compressType,
		resp._nonblock,
		resp._func,
		resp._playload); err != nil {
		return err
	}

	if len(c.Pop()) > 0 {
		goto exit
	}
	if c.Writer().Buffered() > 0 {
		if err = c.Writer().Flush(); err != nil {
			return err
		}
	}
exit:
	return nil
}
