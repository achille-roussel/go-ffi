package ffi

import (
	"strconv"
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

func TestVoidTypeString(t *testing.T) {
	testTypeString(t, Void, "void")
}

func TestCharTypeString(t *testing.T) {
	testTypeString(t, Char, "char")
}

func TestShortTypeString(t *testing.T) {
	testTypeString(t, Short, "short")
}

func TestIntTypeString(t *testing.T) {
	testTypeString(t, Int, "int")
}

func TestLongTypeString(t *testing.T) {
	testTypeString(t, Long, "long")
}

func TestUCharTypeString(t *testing.T) {
	testTypeString(t, UChar, "unsigned char")
}

func TestUShortTypeString(t *testing.T) {
	testTypeString(t, UShort, "unsigned short")
}

func TestUIntTypeString(t *testing.T) {
	testTypeString(t, UInt, "unsigned int")
}

func TestULongTypeString(t *testing.T) {
	testTypeString(t, ULong, "unsigned long")
}

func TestFloatTypeString(t *testing.T) {
	testTypeString(t, Float, "float")
}

func TestDoubleTypeString(t *testing.T) {
	testTypeString(t, Double, "double")
}

func TestInt8TypeString(t *testing.T) {
	testTypeString(t, Int8, "int8_t")
}

func TestInt16TypeString(t *testing.T) {
	testTypeString(t, Int16, "int16_t")
}

func TestInt32TypeString(t *testing.T) {
	testTypeString(t, Int32, "int32_t")
}

func TestInt64TypeString(t *testing.T) {
	testTypeString(t, Int64, "int64_t")
}

func TestUInt8TypeString(t *testing.T) {
	testTypeString(t, UInt8, "uint8_t")
}

func TestUInt16TypeString(t *testing.T) {
	testTypeString(t, UInt16, "uint16_t")
}

func TestUInt32TypeString(t *testing.T) {
	testTypeString(t, UInt32, "uint32_t")
}

func TestUInt64TypeString(t *testing.T) {
	testTypeString(t, UInt64, "uint64_t")
}

func TestPointerTypeString(t *testing.T) {
	testTypeString(t, Pointer, "pointer")
}

func TestDefaultTypeString(t *testing.T) {
	testTypeString(t, Type{}, "struct")
}

func testTypeString(t *testing.T, x Type, s string) {
	if x.String() != s {
		t.Errorf("invalid type string: %s != %s", x, s)
	}
}

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

func TestStatusStringDefault(t *testing.T) {
	if s := Status(-1).String(); s != "unknown" {
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

func TestStatusErrorDefault(t *testing.T) {
	if s := Status(-1).Error(); s != "status: unknown" {
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

	if cif, err = Prepare(Int, Int); err != nil {
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

func TestCallStrerrorReturnString(t *testing.T) {
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

func TestCallStrerrorReturnPointer(t *testing.T) {
	msg := syscall.ENOENT
	ret := unsafe.Pointer(nil)
	arg := int(msg)
	err := Call(unsafe.Pointer(strerror), &ret, arg)

	if err != nil {
		t.Error("call:", err)
		return
	}

	if ret == nil {
		t.Error("call:", ret)
		return
	}
}

func TestCallInvalidReturnTypeNonPointer(t *testing.T) {
	defer func() {
		recover()
	}()

	ret := 0
	arg := -1

	Call(unsafe.Pointer(abs), ret, arg)

	t.Error("unreachable: non-pointer return value should have caused ffi.Call to panic")
}

func TestCallInvalidReturnTypeWrongPointer(t *testing.T) {
	defer func() {
		recover()
	}()

	ret := func() {}
	arg := -1

	Call(unsafe.Pointer(abs), &ret, arg)

	t.Error("unreachable: function pointer return value should have caused ffi.Call to panic")
}

func TestCallInvalidArgumentTypeWrongValue(t *testing.T) {
	defer func() {
		recover()
	}()

	ret := 0
	arg := func() {}

	Call(unsafe.Pointer(abs), &ret, arg)

	t.Error("unreachable: function argument should have caused ffi.Call to panic")
}

func TestCreateAbsClosure(t *testing.T) {
	abs := Closure(func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	})

	if abs == nil {
		t.Error("closure:", abs)
	} else if abs.Pointer() == 0 {
		t.Error("closure: null pointer")
	}
}

func TestCreateItoaClosure(t *testing.T) {
	itoa := Closure(strconv.Itoa)

	if itoa == nil {
		t.Error("closure:", itoa)
	} else if itoa.Pointer() == 0 {
		t.Error("closure: null pointer")
	}
}

func TestCallAbsClosure(t *testing.T) {
	val := 0

	abs := Closure(func(x int) int {
		val = 42

		if x < 0 {
			return -x
		}

		return x
	})

	res := 0
	arg := -1
	err := Call(unsafe.Pointer(abs.Pointer()), &res, arg)

	if err != nil {
		t.Error("closure:", err)
	}

	if res != 1 {
		t.Error("closure: invalid returned value:", res)
	}

	if val != 42 {
		t.Error("closure: invalid value of out-of-scope variable:", val)
	}
}

func TestCallItoaClosure(t *testing.T) {
	itoa := Closure(strconv.Itoa)

	res := ""
	arg := 42
	err := Call(unsafe.Pointer(itoa.Pointer()), &res, arg)

	if err != nil {
		t.Error("closure:", err)
	}

	if res != "42" {
		t.Error("closure: invalid returned value:", res)
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
