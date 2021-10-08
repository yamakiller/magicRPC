package test

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/yamakiller/magicRPC/service"
)

type A struct {
}

func (a *A) TestRpcFuncName(req *service.Request) {

}

func TestFuncName(t *testing.T) {
	testA := &A{}

	//tm := reflect.TypeOf(testA.TestRpcFuncName)

	fmt.Printf("name:%s\n", strings.TrimSuffix(filepath.Ext(runtime.FuncForPC(reflect.ValueOf(testA.TestRpcFuncName).Pointer()).Name()), "-fm"))
	fmt.Printf("name2:%+v\n", reflect.ValueOf(testA.TestRpcFuncName).Pointer())
}
