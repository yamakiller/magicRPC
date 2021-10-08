package service

type Response = Request

func NewReponse(req *Request, ret []byte) *Response {

	resp := &Request{
		_sock:            req._sock,
		_networkProtocol: req._networkProtocol,
		_compressType:    req._compressType,
		_responseStatus:  req._responseStatus,
		_func:            req._func,
		_nonblock:        req._nonblock,
		_timestamp:       req._timestamp,
		_timeout:         req._timeout,
		_sequeNum:        req._sequeNum,
	}

	if ret != nil {
		resp._playload = make([]byte, len(ret))
		copy(resp._playload, ret)
	}

	return resp
}
