package service

import (
	"crypto/tls"
	"errors"
	"reflect"
	"sync"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/sirupsen/logrus"
	"github.com/yamakiller/magicLibs/actors"
	"github.com/yamakiller/magicLibs/boxs"
	"github.com/yamakiller/magicLibs/log"
	"github.com/yamakiller/magicLibs/util"
	"github.com/yamakiller/magicNet/netboxs"
	"github.com/yamakiller/magicNet/netmsgs"
	"github.com/yamakiller/magicRPC/rpc"
	"github.com/yamakiller/magicRPC/rpc/rpcerror"
)

func New(opts ...Option) (*Server, error) {

	opt := defaultOptions
	for _, o := range opts {
		o(&opt)
	}

	if opt._name == "" {
		return nil, rpcerror.ErrNameNeedConfig
	}

	pLogHandle, err := log.SpawnFileLogrus(opt._logLevel, opt._logPath, opt._name)
	if err != nil {
		return nil, err
	}

	pLogAgent := &log.DefaultAgent{}
	pLogAgent.WithHandle(pLogHandle)

	pCore := actors.New(nil)
	pCore.WithLogger(pLogAgent)

	delegates := &delegateBinder{
		_maps: map[uint]*delegateHandle{},
	}

	server := &Server{
		_core:    pCore,
		_clients: make(map[int32]*client),
		_clientPool: &clientPools{
			_core:                pCore,
			_delegateBinderAgent: delegates,
		},
		_clientOfMax:            int32(opt._connMaxOfNumber),
		_clientSendWaitOfNumber: opt._sndWaitOfNumber,
		_flatbufferPool: sync.Pool{
			New: func() interface{} {
				return flatbuffers.NewBuilder(512)
			},
		},
		_delegates: delegates,
		_log:       pLogAgent,
	}

	return server, nil
}

type Options struct {
	_name            string
	_connMaxOfNumber int
	_sndWaitOfNumber int
	_networkProtocol rpc.MRPC_NETWORK_PROTOCOL
	_packageProtocol rpc.MRPC_PACKAGE_PROTOCOL
	_packageCompress rpc.MRPC_PACKAGE_COMPRESS
	_logLevel        logrus.Level
	_logPath         string
}

type Option func(*Options)

var defaultOptions = Options{_connMaxOfNumber: 1024, _sndWaitOfNumber: 8, _logLevel: logrus.DebugLevel, _logPath: "./logs"}

func WithMax(maxOfNumber int) Option {
	return func(opt *Options) {
		opt._connMaxOfNumber = maxOfNumber
	}
}

func WithSendWaitMax(number int) Option {
	return func(opt *Options) {
		opt._sndWaitOfNumber = number
	}
}

func WithName(name string) Option {
	return func(opt *Options) {
		opt._name = name
	}
}

func WithLogLevel(level logrus.Level) Option {
	return func(opt *Options) {
		opt._logLevel = level
	}
}

func WithLogPath(path string) Option {
	return func(opt *Options) {
		opt._logPath = path
	}
}

//Server RPC 服务节点
type Server struct {
	_core                   *actors.Core
	_log                    log.LogAgent
	_network                netboxs.NetBox
	_clients                map[int32]*client
	_clientSync             sync.Mutex
	_clientPool             *clientPools
	_clientOfMax            int32
	_clientSendWaitOfNumber int
	_delegates              *delegateBinder
	_flatbufferPool         sync.Pool
	_networkProtocol        rpc.MRPC_NETWORK_PROTOCOL
	_packageProtocol        rpc.MRPC_PACKAGE_PROTOCOL
	_packageCompress        rpc.MRPC_PACKAGE_COMPRESS
}

func (s *Server) ListenTCPAndServe(addr string) error {
	s._networkProtocol = rpc.MRPC_NETWORK_PROTOCOL_TCP
	s._network = &netboxs.TCPBox{
		Box: *boxs.SpawnBox(nil),
	}
	_, err := s._core.New(func(pid *actors.PID) actors.Actor {
		s._network.(*netboxs.TCPBox).WithPID(pid)
		s._network.(*netboxs.TCPBox).WithMax(int32(s._clientOfMax))
		s._network.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Accept{}), s.onAccept)
		s._network.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Message{}), s.onMessage)
		s._network.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Closed{}), s.onClosed)
		return s._network.(actors.Actor)
	}, actors.PriorityHigh)

	if err != nil {
		return err
	}

	s._network.WithPool(
		&connPools{
			_bfsize: rpc.MRPC_PACKAGE_MAX * 2,
			_wqsize: s._clientSendWaitOfNumber,
		})

	return s._network.ListenAndServe(addr)
}

func (s *Server) ListenTLSAndServe(addr string, ptls *tls.Config) error {
	s._networkProtocol = rpc.MRPC_NETWORK_PROTOCOL_TLS
	s._network = &netboxs.TCPBox{
		Box: *boxs.SpawnBox(nil),
	}

	_, err := s._core.New(func(pid *actors.PID) actors.Actor {
		s._network.(*netboxs.TCPBox).WithPID(pid)
		s._network.(*netboxs.TCPBox).WithMax(int32(s._clientOfMax))
		s._network.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Accept{}), s.onAccept)
		s._network.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Message{}), s.onMessage)
		s._network.(*netboxs.TCPBox).Register(reflect.TypeOf(&netmsgs.Closed{}), s.onClosed)
		return s._network.(actors.Actor)
	}, actors.PriorityHigh)

	if err != nil {
		return err
	}

	return s._network.ListenAndServeTls(addr, ptls)
}

func (s *Server) Shutdown() {
	s._network.Shutdown()
}

func (s *Server) BindDelegate(funcID uint, delegate RequestDelegate, name string) {
	s._delegates.bind(funcID, delegate, name)
}

func (s *Server) MallocFlatBuilder() *flatbuffers.Builder {
	return s._flatbufferPool.Get().(*flatbuffers.Builder)
}

func (s *Server) FreeFlatBuilder(p *flatbuffers.Builder) {
	p.Reset()
	s._flatbufferPool.Put(p)
}

func (s *Server) Response(resp *Response, protocol rpc.MRPC_PACKAGE_PROTOCOL) error {
	if resp.getCompressType() == rpc.MRPC_CT_NO {
		resp.withCompressType(s._packageCompress)
	}

	return s._network.SendTo(resp._sock, resp)
}

func (s *Server) onAccept(context *boxs.Context) {
	request := context.Message().(*netmsgs.Accept)
	if request == nil {
		context.Error("accept message error")
		return
	}

	if err := s.registerClient(request.Sock); err != nil {
		context.Error("accept register error:%v", err)
		return
	}
	s._network.OpenTo(request.Sock)
}

func (s *Server) onMessage(context *boxs.Context) {
	request := context.Message().(*netmsgs.Message)
	if request == nil {
		context.Error("message error")
		return
	}

	pk := request.Data.(*rpc.Packet)
	util.AssertEmpty(pk, "network data packet is null")

	switch pk.GetProtocolType() {
	case rpc.MRPC_PT_JSON:
	case rpc.MRPC_PT_XML:
	case rpc.MRPC_PT_PROTOBUFF:
	case rpc.MRPC_PT_FLATBUFF:
	default:
		context.Warning("unpack serialized type is failed!, type:%d", pk.GetProtocolType())
		return
	}

	rpcReq := &Request{}
	rpcReq.bindSock(request.Sock)
	rpcReq.withNetworkProtocol(s._networkProtocol)
	rpcReq.withSequence(pk.Serialized())
	rpcReq.withNonblock(pk.IsNonblok())
	rpcReq.withCompressType(pk.GetCompressType())
	rpcReq.bindFunc(pk.Func())
	rpcReq.Push(pk.Payload())
	rpcReq.WithStatus(rpc.RS_OK)
	if pk.GetStatusCode() != rpc.MRPSC_OK {
		rpcReq.WithStatus(rpc.RS_FAILD)
	}

	client := s.getClient(request.Sock.(int32))
	if client == nil {
		context.Error("socket not found:%d", request.Sock)
		return
	}
	defer s.freeClient(client)

	client.GetPID().Post(rpcReq)

	context.Debug("request event %d", request.Sock)
}

func (s *Server) onClosed(context *boxs.Context) {
	request := context.Message().(*netmsgs.Closed)
	if request == nil {
		context.Error("socket close error")
		return
	}

	s.removeClient(request.Sock)
	context.Debug("close socket %d", request.Sock)
}

func (s *Server) getClient(sock int32) *client {
	s._clientSync.Lock()
	defer s._clientSync.Unlock()

	if _, ok := s._clients[sock]; !ok {
		return nil
	}

	r := s._clients[sock]
	r.refInc()
	return r
}

func (s *Server) freeClient(c *client) {
	s._clientSync.Lock()
	ref := c.refDec()
	if ref <= 0 {
		delete(s._clients, c._socket)
		s._clientSync.Unlock()
		s._clientPool.put(c)
		return
	}
	s._clientSync.Unlock()

}

func (s *Server) registerClient(sock int32) error {
	s._clientSync.Lock()
	defer s._clientSync.Unlock()

	if len(s._clients) > 0 {
		if _, ok := s._clients[sock]; ok {
			return errors.New("Socket already exists")
		}
	}
	p := s._clientPool.get()
	p._socket = sock
	p._ref = 1
	s._clients[sock] = p

	return nil
}

func (s *Server) removeClient(sock int32) {
	s._clientSync.Lock()
	p := s._clients[sock]
	if p == nil {
		s._clientSync.Unlock()
		return
	}

	delete(s._clients, sock)

	ref := p.refDec()
	s._clientSync.Unlock()

	if ref <= 0 {
		s._clientPool.put(p)
	}
}
