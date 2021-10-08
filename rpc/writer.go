package rpc

import "io"

func writes(w io.Writer, data []byte) error {
	var (
		offset, count, ret int
		err                error
	)

	count = len(data)
	for {
		ret, err = w.Write(data[offset:])
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
