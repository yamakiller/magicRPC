package client

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/log"
	magiclibconn "github.com/yamakiller/magicLibs/net/connection"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicRPC/rpc"
	"github.com/yamakiller/magicRPC/rpc/rpcerror"
	"github.com/yamakiller/magicRPC/rpc/rpctime"
)

func New(opts ...Option) (*Broker, error) {
	opt := defaultOptions
	for _, o := range opts {
		o(&opt)
	}

	if opt._name == "" {
		return nil, rpcerror.ErrNameNeedConfig
	}

	if opt._logAgent == nil {
		return nil, rpcerror.ErrLoggerNeedConfig
	}

	broker := &Broker{
		_name:              opt._name,
		_clientRequestSn:   0,
		_clientRequestSync: 0,
		_clientRecv:        make(chan *Request),
		_packageProtocol:   opt._packageProtocol,
		_packageCompress:   opt._packageCompress,
		_log:               opt._logAgent,
	}

	broker._client = &magiclibconn.TCPClient{
		ReadBufferSize:  rpc.MRPC_PACKAGE_MAX * 2,
		WriteBufferSize: rpc.MRPC_PACKAGE_MAX * 2,
		WriteWaitQueue:  opt._sendBufferOfNumber,
		S:               broker,
		E:               broker,
	}

	return broker, nil
}

type Options struct {
	_name               string
	_sendBufferOfNumber int
	_logAgent           log.LogAgent
	_packageProtocol    rpc.MRPC_PACKAGE_PROTOCOL
	_packageCompress    rpc.MRPC_PACKAGE_COMPRESS
}

type Option func(*Options)

func WithName(name string) Option {
	return func(opt *Options) {
		opt._name = name
	}
}

func WithSendBufferOfNumber(number int) Option {
	return func(opt *Options) {
		opt._sendBufferOfNumber = number
	}
}

func WithLogAgent(plog log.LogAgent) Option {
	return func(opt *Options) {
		opt._logAgent = plog
	}
}

func WithPackageProtocol(proto rpc.MRPC_PACKAGE_PROTOCOL) Option {
	return func(opt *Options) {
		opt._packageProtocol = proto
	}
}

func WithPackageCompress(compress rpc.MRPC_PACKAGE_COMPRESS) Option {
	return func(opt *Options) {
		opt._packageCompress = compress
	}
}

var defaultOptions = Options{_sendBufferOfNumber: 8,

	_packageProtocol: rpc.MRPC_PT_FLATBUFF,
	_packageCompress: rpc.MRPC_CT_NO,
}

type Broker struct {
	_name              string
	_client            magiclibconn.Client
	_clientRequestSn   uint32
	_clientRequestSync uint32
	_clientRecv        chan *Response
	_log               log.LogAgent
	_networkProtocol   rpc.MRPC_NETWORK_PROTOCOL
	_packageProtocol   rpc.MRPC_PACKAGE_PROTOCOL
	_packageCompress   rpc.MRPC_PACKAGE_COMPRESS
}

func (b *Broker) Connect(addr string, timeout int) error {
	b._networkProtocol = rpc.MRPC_NETWORK_PROTOCOL_TCP
	if err := b._client.Connect(addr, time.Duration(timeout)*time.Millisecond); err != nil {
		return err
	}

	go b.recvServe()
	return nil
}

func (b *Broker) TLSConnect(addr string, config *tls.Config, timeout int) error {
	b._networkProtocol = rpc.MRPC_NETWORK_PROTOCOL_TLS
	if err := b._client.ConnectTls(addr, time.Duration(timeout)*time.Millisecond, config); err != nil {
		return err
	}
	return nil
}

func (b *Broker) SyncCall(req *Request, prootocol rpc.MRPC_PACKAGE_PROTOCOL) (*Response, error) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(req._timeout)*time.Millisecond)
	defer cancel()

	sn := b.getSerialOfNumber()
	req._sequeNum = sn
	req._timestamp = rpctime.Click()

	if err := b._client.SendTo(req); err != nil {
		b._clientRequestSync = 0
		return nil, err
	}

	select {
	case resp := <-b._clientRecv:
		return resp, nil
	case <-b._client.IsConnected().(context.Context).Done():
		return nil, errors.New("connection is  disconnect")
	case <-ctx.Done():
		return nil, rpcerror.ErrBrokerRequestTimeOut
	}
}

func (b *Broker) Close() {
	b._client.Close()

	close(b._clientRecv)
}

func (b *Broker) Error(err error) {
	b._log.Error(b._name, err.Error())
}

func (b *Broker) ErrorLog(fmt string, args ...interface{}) {
	b._log.Error(b._name, fmt, args...)
}

func (b *Broker) InfoLog(fmt string, args ...interface{}) {
	b._log.Info(b._name, fmt, args...)
}

func (b *Broker) DebugLog(fmt string, args ...interface{}) {
	b._log.Debug(b._name, fmt, args...)
}

func (b *Broker) UnSeria(r io.Reader) (interface{}, int, error) {
	pk, err := rpc.Decoding(r)
	if err != nil {
		return nil, -1, err
	}

	//kleepalive data
	if pk.Func() == 0 {
		//TODO: 心跳回复
		return nil, 0, nil
	}

	return pk, pk.Size(), nil
}

func (b *Broker) Seria(msg interface{}, w io.Writer) (int, error) {
	req := msg.(*Request)
	util.AssertEmpty(req, b._name+" rpc broker seria fail, not *client.Request")
	var (
		err error
		ret int
	)
	if ret, err = rpc.Encoding(w,
		b._packageProtocol,
		req._sequeNum,
		req._compressType,
		req._nonblock,
		req._func,
		req._payload); err != nil {
		return -1, err
	}

	return ret, nil
}

func (b *Broker) recvServe() {
	for {
		pk, err := b._client.Parse()
		if err != nil {
			return
		}
		b.onMessage(pk)
	}
}

func (b *Broker) onMessage(msg interface{}) {
	pk := msg.(*rpc.Packet)
	util.AssertEmpty(pk, fmt.Sprintf("%s rpc message to packet fail, %v", b._name, msg))
	resp := &Response{
		_compressType:   pk.GetCompressType(),
		_responseStatus: rpc.RESPONSE_STATUS(pk.GetStatusCode()),
		_func:           pk.Func(),
		_nonblock:       pk.IsNonblok(),
		_sequeNum:       pk.Serialized(),
	}

	if pk.Payload() != nil {
		payload := make([]byte, len(pk.Payload()))
		copy(payload, pk.Payload())
		switch resp._compressType {
		case rpc.MRPC_CT_DYNAIC:
		case rpc.MRPC_CT_COMPRESS:
		default:
		}
		resp.Push(payload)
	}

	if resp._sequeNum == 0 && pk.Func() == 0 {
		//心跳
		keepReq := &Request{
			_compressType:   b._packageCompress,
			_responseStatus: rpc.RS_OK,
			_func:           0,
			_nonblock:       false,
			_timestamp:      0,
			_timeout:        0,
			_sequeNum:       0,
		}

		if err := b._client.SendTo(keepReq); err != nil {
			b.ErrorLog(b._name, "keepalive response fail:%v", err)
			return
		}

		return
	}

	if atomic.CompareAndSwapUint32(&b._clientRequestSync, resp._sequeNum, 0) {
		b._clientRecv <- resp
		return
	}

	b.ErrorLog("response undefine register task %d", resp._sequeNum)
}

func (b *Broker) getSerialOfNumber() uint32 {
	for {
		sn := atomic.AddUint32(&b._clientRequestSn, 1)
		if sn != 0 {
			return sn
		}
	}
}
