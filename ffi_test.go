package ffi

import (
	"testing"
	"unsafe"

	"github.com/achille-roussel/go-dl"
)

var (
	libc dl.Library
	abs  uintptr
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

	var arg int32 = -1
	var res int32

	if err = cif.Call(unsafe.Pointer(abs), unsafe.Pointer(&res), unsafe.Pointer(&arg)); err != nil {
		t.Error("call:", err)
		return
	}

	if res != 1 {
		t.Error("call:", res)
		return
	}
}

func init() {
	var err error

	if libc, err = load("libc"); err != nil {
		panic(err)
	}

	if abs, err = libc.Symbol("abs"); err != nil {
		panic(err)
	}
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
