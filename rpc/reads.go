package rpc

import "io"

func reads(r io.Reader, out []byte) error {
	var (
		offset, count, ret int
		err                error
	)

	count = len(out)
	for {
		ret, err = r.Read(out[offset:])
		if err != nil {
			return err
		}
		offset += ret
		if offset >= count {
			break
		}
	}

	return nil
}
