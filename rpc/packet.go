package rpc

import (
	"io"

	"github.com/yamakiller/magicRPC/rpc/rpcerror"
)

type Packet struct {
	_header  Header
	_payload []byte
}

func (p *Packet) Size() int {
	return p._header.GetPackageSize() + p._header.GetBodySize()
}

func (p *Packet) Func() uint32 {
	return p._header._func
}

func (p *Packet) Serialized() uint32 {
	return p._header._sequence
}

func (p *Packet) IsNonblok() bool {
	return p._header.getNonblock()
}

func (p *Packet) GetProtocolType() MRPC_PACKAGE_PROTOCOL {
	return p._header.getProtocolType()
}

func (p *Packet) GetCompressType() MRPC_PACKAGE_COMPRESS {
	return p._header.getCompressType()
}

func (p *Packet) GetStatusCode() MRPC_PACKAGE_STATUS_CODE {
	return p._header.getStatusCode()
}

func (p *Packet) Payload() []byte {
	return p._payload
}

func (p *Packet) unmarshal(r io.Reader) error {
	if err := p._header.unmarshal(r); err != nil {
		return err
	}

	if p._header._bodysize == 0 && p._header._func != 0 && p._header._sequence != 0 {
		return rpcerror.ErrProtocolKeepalive
	}

	if p._header.GetBodySize()+p._header.GetPackageSize() > MRPC_PACKAGE_MAX {
		return rpcerror.ErrProtocolDataOverflow
	}

	if p._header._bodysize > 0 {
		p._payload = make([]byte, p._header._bodysize)
		if err := reads(r, p._payload); err != nil {
			return err
		}

		if p._header.isCompress() {
			//TODO: 解压
		}
	}

	return nil
}
