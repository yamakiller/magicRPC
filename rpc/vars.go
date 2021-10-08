package rpc

var (
	_MRPC_PACKAGE_SINGLE_MAX = 1024 * 64
)

func WithPackageSingleMax(max int) {
	_MRPC_PACKAGE_SINGLE_MAX = max
}

func GetPackageSingleMax() int {
	return _MRPC_PACKAGE_SINGLE_MAX
}
