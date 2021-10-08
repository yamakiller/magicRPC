package rpc

import (
	"io"

	"github.com/yamakiller/magicRPC/rpc/rpcerror"
)

func Encoding(w io.Writer,
	tpp MRPC_PACKAGE_PROTOCOL,
	sequence uint32,
	compressType MRPC_PACKAGE_COMPRESS,
	isNonblock bool,
	fun uint,
	payload []byte) (int, error) {
	h := Header{}
	h.Init(tpp, sequence, compressType, isNonblock, fun)

	if payload != nil {
		switch h.getCompressType() {
		case MRPC_CT_COMPRESS:
			return -1, rpcerror.ErrUndefineCompressMode
		case MRPC_CT_DYNAIC:
			return -1, rpcerror.ErrUndefineCompressMode
		default:
		}
	}

	h._bodysize = 0
	if payload != nil {
		h._bodysize = uint16(len(payload))
	}
	packageSize := h.GetPackageSize()

	if packageSize+h.GetBodySize() > MRPC_PACKAGE_MAX {
		return -1, rpcerror.ErrProtocolDataOverflow
	}

	var (
		err error
	)

	if err = h.marshal(w); err != nil {
		return -1, err
	}

	if payload != nil {
		if err = writes(w, payload); err != nil {
			return -1, err
		}
	}

	return packageSize + h.GetBodySize(), nil
}
