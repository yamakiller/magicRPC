package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/yamakiller/magicLibs/log"
	mlibconn "github.com/yamakiller/magicLibs/net/connection"
	mlibtcp "github.com/yamakiller/magicLibs/net/connection/tcp"
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

	broker._client = mlibtcp.New(rpc.MRPC_PACKAGE_MAX*2, opt._sendBufferOfNumber)
	broker._client.(*mlibtcp.Client).Serializer = broker
	broker._client.(*mlibtcp.Client).Exception = broker

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
	_client            mlibconn.Client
	_clientRequestSn   uint32
	_clientRequestSync uint32
	_clientRecv        chan *Response
	_clientWait        sync.WaitGroup
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
	b._clientWait.Add(1)
	defer b._clientWait.Done()

	if b._clientRecv == nil {
		return nil, errors.New("broker connect destoryed")
	}

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
	case <-ctx.Done():
		return nil, rpcerror.ErrBrokerRequestTimeOut
	}
}

func (b *Broker) Close() {
	b._client.Close()
}

func (b *Broker) Shutdown() {
	if b._clientRecv != nil {
		close(b._clientRecv)
	}
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
	if err != nil && err == rpcerror.ErrIsProtocolKeepalive {
		//TODO: 心跳回复
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
			return nil, -1, err
		}

		return nil, 0, nil
	}

	if err != nil {
		return nil, -1, err
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
	h := rpc.Header{}
	h.Init(b._packageProtocol, req._sequeNum, req._compressType, req._nonblock, req._func)
	packageSize := h.GetPackageSize()
	if len(req._payload) > (rpc.MRPC_PACKAGE_MAX*2 - w.(*bufio.Writer).Buffered() - packageSize) {
		if w.(*bufio.Writer).Buffered() > 0 {
			if err = w.(*bufio.Writer).Flush(); err != nil {
				return ret, err
			}
		}
	}
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
		if pk == nil {
			continue
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
