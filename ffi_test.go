package ffi

import (
	"strings"
	"syscall"
	"testing"
	"unsafe"

	"github.com/achille-roussel/go-dl"
)

var (
	libc     dl.Library
	abs      uintptr
	strerror uintptr
)

func TestStatusStringOK(t *testing.T) {
	if s := OK.String(); s != "OK" {
		t.Error("invalid status string:", s)
	}
}

func TestStatusStringBadTypedef(t *testing.T) {
	if s := BadTypedef.String(); s != "bad-typedef" {
		t.Error("invalid status string:", s)
	}
}

func TestStatusStringBadABI(t *testing.T) {
	if s := BadABI.String(); s != "bad-ABI" {
		t.Error("invalid status string:", s)
	}
}

func TestStatusErrorOK(t *testing.T) {
	if s := OK.Error(); s != "status: OK" {
		t.Error("invalid status string:", s)
	}
}

func TestStatusErrorBadTypedef(t *testing.T) {
	if s := BadTypedef.Error(); s != "status: bad-typedef" {
		t.Error("invalid status string:", s)
	}
}

func TestStatusErrorBadABI(t *testing.T) {
	if s := BadABI.Error(); s != "status: bad-ABI" {
		t.Error("invalid status string:", s)
	}
}

func TestPrepareVoid(t *testing.T) {
	if _, err := Prepare(Void); err != nil {
		t.Error(err)
	}
}

func TestPrepareAndCall(t *testing.T) {
	var cif Interface
	var err error

	if cif, err = Prepare(Sint, Sint); err != nil {
		t.Error("prepare:", err)
		return
	}

	ret := 0
	arg := -1

	if err = cif.Call(unsafe.Pointer(abs), unsafe.Pointer(&ret), unsafe.Pointer(&arg)); err != nil {
		t.Error("call:", err)
		return
	}

	if ret != 1 {
		t.Error("call:", ret)
		return
	}
}

func TestCallAbs(t *testing.T) {
	ret := 0
	arg := -1
	err := Call(unsafe.Pointer(abs), &ret, arg)

	if err != nil {
		t.Error("call:", err)
		return
	}

	if ret != 1 {
		t.Error("call:", ret)
		return
	}
}

func TestCallStrerror(t *testing.T) {
	msg := syscall.ENOENT
	ret := ""
	arg := int(msg)
	err := Call(unsafe.Pointer(strerror), &ret, arg)

	if err != nil {
		t.Error("call:", err)
		return
	}

	if strings.ToLower(ret) != msg.Error() {
		t.Error("call:", ret)
		return
	}
}

func init() {
	var err error

	if libc, err = load("libc"); err != nil {
		panic(err)
	}

	abs = symbol(libc, "abs")
	strerror = symbol(libc, "strerror")
}

func load(name string) (lib dl.Library, err error) {
	var path string

	if path, err = dl.Find(name); err != nil {
		return
	}

	if lib, err = dl.Open(path, 0); err != nil {
		return
	}

	return
}

func symbol(lib dl.Library, name string) (addr uintptr) {
	var err error

	if addr, err = libc.Symbol(name); err != nil {
		panic(err)
	}

	return
}
