package ffi

import (
	"fmt"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"unsafe"

	"github.com/achille-roussel/go-dl"
)

var (
	libc     dl.Library
	libm     dl.Library
	abs      uintptr
	fabs     uintptr
	fabsf    uintptr
	snprintf uintptr
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
	testTypeString(t, Pointer, "void *")
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
	Prepare(Void)
}

func TestPrepareAndCall(t *testing.T) {
	cif := Prepare(Int, Int)
	ret := 0
	arg := -1

	if err := cif.Call(unsafe.Pointer(abs), unsafe.Pointer(&ret), unsafe.Pointer(&arg)); err != nil {
		t.Error("call:", err)
		return
	}

	if ret != 1 {
		t.Error("call:", ret)
		return
	}
}

func TestInterfaceString(t *testing.T) {
	if s := Prepare(Void, Pointer, Int).String(); s != "void(*)(void *, int)" {
		t.Error("invalid string representation of call interface:", s)
	}
}

func TestCallAbs(t *testing.T) {
	ret := 0
	arg := -1
	err := Call(unsafe.Pointer(abs), &ret, arg)

	if err != nil {
		t.Error("abs:", err)
		return
	}

	if ret != 1 {
		t.Error("abs:", ret)
		return
	}
}

func TestCallFabs(t *testing.T) {
	res := float64(0.0)
	arg := float64(-0.5)
	err := Call(unsafe.Pointer(fabs), &res, arg)

	if err != nil {
		t.Error("fabs:", err)
	}

	if res != 0.5 {
		t.Error("fabs:", res)
	}
}

func TestCallFabsf(t *testing.T) {
	res := float32(0.0)
	arg := float32(-0.5)
	err := Call(unsafe.Pointer(fabsf), &res, arg)

	if err != nil {
		t.Error("fabsf:", err)
	}

	if res != 0.5 {
		t.Error("fabsf:", res)
	}
}

func TestCallStrerrorReturnString(t *testing.T) {
	msg := syscall.ENOENT
	ret := ""
	arg := int(msg)
	err := Call(unsafe.Pointer(strerror), &ret, arg)

	if err != nil {
		t.Error("strerror:", err)
		return
	}

	if strings.ToLower(ret) != msg.Error() {
		t.Error("strerror:", ret)
		return
	}
}

func TestCallStrerrorReturnPointer(t *testing.T) {
	msg := syscall.ENOENT
	ret := unsafe.Pointer(nil)
	arg := int(msg)
	err := Call(unsafe.Pointer(strerror), &ret, arg)

	if err != nil {
		t.Error("strerror:", err)
		return
	}

	if ret == nil {
		t.Error("strerror:", ret)
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

func TestCallSnprintfInt(t *testing.T) {
	testCallSnprintf(t, "%d", int(42))
}

func TestCallSnprintfInt8(t *testing.T) {
	testCallSnprintf(t, "%d", int8(42))
}

func TestCallSnprintfInt16(t *testing.T) {
	testCallSnprintf(t, "%d", int16(42))
}

func TestCallSnprintfInt32(t *testing.T) {
	testCallSnprintf(t, "%d", int32(42))
}

func TestCallSnprintfInt64(t *testing.T) {
	testCallSnprintf(t, "%ld", int64(42))
}

func TestCallSnprintfUint(t *testing.T) {
	testCallSnprintf(t, "%u", uint(42))
}

func TestCallSnprintfUint8(t *testing.T) {
	testCallSnprintf(t, "%u", uint8(42))
}

func TestCallSnprintfUint16(t *testing.T) {
	testCallSnprintf(t, "%u", uint16(42))
}

func TestCallSnprintfUint32(t *testing.T) {
	testCallSnprintf(t, "%u", uint32(42))
}

func TestCallSnprintfUint64(t *testing.T) {
	testCallSnprintf(t, "%lu", uint64(42))
}

func TestCallSnprintfFloat64(t *testing.T) {
	testCallSnprintf(t, "%g", float64(42))
}

func TestCallSnprintfString(t *testing.T) {
	testCallSnprintf(t, "%s", "Hello World!")
}

func testCallSnprintf(t *testing.T, f string, v interface{}) {
	buf := make([]byte, 128)
	res := 0
	err := Call(unsafe.Pointer(snprintf), &res, &buf[0], uintptr(len(buf)), f, v)
	ref := fmt.Sprint(v)

	if err != nil {
		t.Error("snprintf:", err)
	}

	if res != len(ref) {
		t.Error("snprintf: invalid return value:", res, "!=", len(ref))
	}

	if s := string(buf[:res]); s != ref {
		t.Error("snprintf: invalid formatted string:", s, "!=", ref)
	}
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

	if libm, err = load("libm"); err != nil {
		panic(err)
	}

	abs = symbol(libc, "abs")
	fabs = symbol(libm, "fabs")
	fabsf = symbol(libm, "fabsf")
	snprintf = symbol(libc, "snprintf")
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

	if addr, err = lib.Symbol(name); err != nil {
		panic(err)
	}

	return
}

func BenchmarkCallingAbsViaCgo(b *testing.B) {
	for i, n := 0, b.N; i != n; i++ {
		ffi_test_abs__(-i)
	}
}

func BenchmarkCallingAbsViaInterface(b *testing.B) {
	cif := Prepare(Int, Int)

	for i, n := 0, b.N; i != n; i++ {
		arg := -i
		res := 0
		cif.Call(unsafe.Pointer(abs), unsafe.Pointer(&res), unsafe.Pointer(&arg))
	}
}

func BenchmarkCallingAbsViaCall(b *testing.B) {
	for i, n := 0, b.N; i != n; i++ {
		arg := -1
		res := 0
		Call(unsafe.Pointer(abs), &res, arg)
	}
}

func BenchmarkCallingAbsViaClosure(b *testing.B) {
	abs := Closure(func(n int) int {
		if n < 0 {
			return -n
		}
		return n
	})

	for i, n := 0, b.N; i != n; i++ {
		arg := -i
		res := 0
		abs.Call(unsafe.Pointer(&res), unsafe.Pointer(&arg))
	}
}
