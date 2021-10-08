package rpctime

import "time"

var (
	_offsetTimespace int64
)

func init() {
	utc, _ := time.LoadLocation("")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2021-01-01 00:00:00", utc)
	_offsetTimespace = t.UnixNano() / 1e6 //time.UTC().Format("2021-01-01 00:00:00").UnixNano() / 1e6
}

func Click() uint32 {
	timespace := time.Now().UTC().UnixNano() / 1e6
	return uint32(timespace - _offsetTimespace)
}
