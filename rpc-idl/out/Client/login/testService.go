//magic RPC Client Automatic generation service:2021-10-07 11:37:42

package login

import (
	rpclient "github.com/yamakiller/magicRPC/client"
	rpcc "github.com/yamakiller/magicRPC/rpc"
)


type TestService struct {
	_brokerImpl *rpclient.Broker
	_compressType rpcc.MRPC_PACKAGE_COMPRESS
}

func (t *TestService) WithBroker(broker *rpclient.Broker) {
	t._brokerImpl = broker
}
func (t *TestService) WithCommpressType(compressType rpcc.MRPC_PACKAGE_COMPRESS) {
	t._compressType = compressType
}
func (t *TestService) GetName(sign Signin, timeout int) (*UserInfo, error) {
	return testInterface{}.getNameInterface(t._brokerImpl, sign, t._compressType, timeout)
}

