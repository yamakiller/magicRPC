package rpc

import "io"

func Decoding(r io.Reader) (*Packet, error) {
	pk := Packet{}
	if err := pk.unmarshal(r); err != nil {
		return nil, err
	}

	return &pk, nil
}
