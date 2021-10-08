package service

import (
	"github.com/yamakiller/magicLibs/boxs"
)

type RequestDelegate func(*boxs.Context, *Request) error
type delegateGet func(uint) *delegateHandle

type delegateHandle struct {
	_func RequestDelegate
	_name string
}

type delegateBinder struct {
	_maps map[uint]*delegateHandle
}

func (d *delegateBinder) bind(id uint, delegate RequestDelegate, name string) {
	//strings.TrimSuffix(filepath.Ext(runtime.FuncForPC(reflect.ValueOf(delegate).Pointer()).Name()), "-fm")
	d._maps[id] = &delegateHandle{
		_func: delegate,
		_name: name,
	}
}

func (d *delegateBinder) get(id uint) *delegateHandle {
	if _, ok := d._maps[id]; !ok {
		return nil
	}

	return d._maps[id]
}
