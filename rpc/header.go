package rpc

import (
	"encoding/binary"
	"io"
	"unsafe"

	"github.com/yamakiller/magicRPC/rpc/rpcerror"
)

type Header struct {
	_flag     []byte
	_ver      uint8
	_st       uint16
	_status   uint8
	_reserved [2]byte
	_sequence uint32
	_func     uint
	_bodysize uint16
}

func (h *Header) Init(tpp MRPC_PACKAGE_PROTOCOL, sequence uint32, compressType MRPC_PACKAGE_COMPRESS, isNonblock bool, fun uint) {
	h._flag = _MRPC_HEADER_FLAG[:]
	h._ver = _RPC_VERSION
	h._sequence = sequence
	h._func = fun
	h._st = uint16(tpp) & _MRPC_HEADER_PROTOCOL_MASK
	h._st |= uint16((uint16(compressType) << _MRPC_HANDER_STATUS_CODE_SHIFT))
	h._status = 0
	if isNonblock {
		h._status |= _MRPC_HANDER_NONBLOCK_MASK
	}
	h._status |= (uint8(MRPSC_OK) & _MRPC_HANDER_STATUS_CODE_MASK)
}

func (h *Header) GetPackageSize() int {
	return len(h._flag) + h.size() + _MRPC_PREFIX_SIZE
}

func (h *Header) GetFunc() uint {
	return h._func
}

func (h *Header) GetVersion() uint8 {
	return h._ver
}

func (h *Header) GetBodySize() int {
	return int(h._bodysize)
}

func (h *Header) size() int {
	return int(unsafe.Sizeof(h._ver)) + int(unsafe.Sizeof(h._st)) + int(unsafe.Sizeof(h._status)) + int(unsafe.Sizeof(h._sequence)) +
		len(h._reserved) + int(unsafe.Sizeof(h._func)) + int(unsafe.Sizeof(h._bodysize))
}

func (h *Header) getProtocolType() MRPC_PACKAGE_PROTOCOL {
	return MRPC_PACKAGE_PROTOCOL(h._st & _MRPC_HEADER_PROTOCOL_MASK)
}

func (h *Header) getCompressType() MRPC_PACKAGE_COMPRESS {
	return MRPC_PACKAGE_COMPRESS((h._st & _MRPC_HANDER_CPNPRESS_TYPE_MASK) >> _NRPC_HEADER_CPNPRESS_TYPE_SHIFT)
}

func (h *Header) getSerialized() uint16 {
	return uint16(h._sequence)
}

func (h *Header) getNonblock() bool {
	if (h._status & _MRPC_HANDER_NONBLOCK_MASK) != 0 {
		return false
	}
	return true
}

func (h *Header) getStatusCode() MRPC_PACKAGE_STATUS_CODE {
	return MRPC_PACKAGE_STATUS_CODE(uint8(h._status) & uint8(_MRPC_HANDER_STATUS_CODE_MASK))
}

func (h *Header) WithStatusCode(status MRPC_PACKAGE_STATUS_CODE) {
	h._status |= (uint8(status) & _MRPC_HANDER_STATUS_CODE_MASK)
}

func (h *Header) getHB() bool {
	if (h._status & _MRPC_HANDER_HB_MASK) != 0 {
		return true
	}
	return false
}

func (h *Header) WithHBPackage() {
	h._status |= _MRPC_HANDER_HB_MASK
}

func (h *Header) getOW() bool {
	if (h._status & _MRPC_HANDER_OW_MASK) != 0 {
		return true
	}
	return false
}

func (h *Header) WithOWPackage() {
	h._status |= _MRPC_HANDER_OW_MASK
}

func (h *Header) getRP() bool {
	if (h._status & _MRPC_HANDER_RP_MASK) != 0 {
		return true
	}
	return false
}

func (h *Header) WithRPPackage() {
	h._status |= _MRPC_HANDER_RP_MASK
}

func (h *Header) isCompress() bool {
	if h._st&_MRPC_HANDER_CPNPRESS_TYPE_MASK != 0 {
		return true
	}
	return false
}

func (h *Header) marshal(w io.Writer) error {
	var (
		err error
	)
	if err = writes(w, h._flag[:]); err != nil {
		return err
	}

	if err = binary.Write(w, binary.BigEndian, uint8(h.size())); err != nil {
		return err
	}

	if err = binary.Write(w, binary.BigEndian, h._ver); err != nil {
		return err
	}

	if err = binary.Write(w, binary.BigEndian, h._st); err != nil {
		return err
	}

	if err = binary.Write(w, binary.BigEndian, h._status); err != nil {
		return err
	}

	if err = writes(w, h._reserved[:]); err != nil {
		return err
	}

	if getSys() == _SYS64 {
		_func64 := uint64(h._func)
		if err = binary.Write(w, binary.BigEndian, _func64); err != nil {
			return err
		}
	} else {
		_func32 := uint32(h._func)
		if err = binary.Write(w, binary.BigEndian, _func32); err != nil {
			return err
		}
	}

	if err = binary.Write(w, binary.BigEndian, h._sequence); err != nil {
		return err
	}

	if err = binary.Write(w, binary.BigEndian, h._bodysize); err != nil {
		return err
	}
	return nil
}

func (h *Header) unmarshal(r io.Reader) error {
	var (
		flag  [3]byte
		hsize uint8
		err   error
	)
	if err = reads(r, flag[:]); err != nil {
		return err
	}

	if flag[0] != _MRPC_HEADER_FLAG[0] || flag[1] != _MRPC_HEADER_FLAG[1] || flag[2] != _MRPC_HEADER_FLAG[2] {
		return rpcerror.ErrProtocol
	}

	if err = binary.Read(r, binary.BigEndian, &hsize); err != nil {
		return err
	}

	if err = binary.Read(r, binary.BigEndian, &h._ver); err != nil {
		return err
	}

	if h._ver != _RPC_VERSION {
		return rpcerror.ErrProtocolVersion
	}

	if err = binary.Read(r, binary.BigEndian, &h._st); err != nil {
		return err
	}

	if err = binary.Read(r, binary.BigEndian, &h._status); err != nil {
		return err
	}

	if err = reads(r, h._reserved[:]); err != nil {
		return err
	}
	if getSys() == _SYS64 {
		_func64 := uint64(h._func)
		if err = binary.Read(r, binary.BigEndian, &_func64); err != nil {
			return err
		}
		h._func = uint(_func64)
	} else {
		_func32 := uint32(h._func)
		if err = binary.Read(r, binary.BigEndian, &_func32); err != nil {
			return err
		}
		h._func = uint(_func32)
	}

	if err = binary.Read(r, binary.BigEndian, &h._sequence); err != nil {
		return err
	}

	if err = binary.Read(r, binary.BigEndian, &h._bodysize); err != nil {
		return err
	}
	//检测协议是否符号标准
	if hsize != uint8(h.size()) {
		return rpcerror.ErrProtocolHeaderValid
	}

	if h.GetPackageSize() > _MRPC_PACKAGE_SINGLE_MAX {
		return rpcerror.ErrProtocolDataOverflow
	}

	return nil
}
