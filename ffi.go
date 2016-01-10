package ffi

// #include <ffi.h>
// #include <stdint.h>
//
// typedef void (*function)(void);
//
import "C"
import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"
)

type Status int

const (
	OK         Status = Status(C.FFI_OK)
	BadTypedef Status = Status(C.FFI_BAD_TYPEDEF)
	BadABI     Status = Status(C.FFI_BAD_ABI)
)

func (s Status) String() string {
	switch s {
	case OK:
		return "OK"
	case BadTypedef:
		return "bad-typedef"
	case BadABI:
		return "bad-ABI"
	default:
		return "unknown"
	}
}

func (s Status) Error() string {
	return "status: " + s.String()
}

type Type struct {
	ffi_type *C.ffi_type
	name     string
}

var (
	Void Type = Type{&C.ffi_type_void, "void"}

	UChar  Type = Type{&C.ffi_type_uchar, "unsigned char"}
	UShort Type = Type{&C.ffi_type_ushort, "unsigned short"}
	UInt   Type = Type{&C.ffi_type_uint, "unsigned int"}
	ULong  Type = Type{&C.ffi_type_ulong, "unsigned long"}

	UInt8  Type = Type{&C.ffi_type_uint8, "uint8_t"}
	UInt16 Type = Type{&C.ffi_type_uint16, "uint16_t"}
	UInt32 Type = Type{&C.ffi_type_uint32, "uint32_t"}
	UInt64 Type = Type{&C.ffi_type_uint64, "uint64_t"}

	Char  Type = Type{&C.ffi_type_schar, "char"}
	Short Type = Type{&C.ffi_type_sshort, "short"}
	Int   Type = Type{&C.ffi_type_sint, "int"}
	Long  Type = Type{&C.ffi_type_slong, "long"}

	Int8  Type = Type{&C.ffi_type_sint8, "int8_t"}
	Int16 Type = Type{&C.ffi_type_sint16, "int16_t"}
	Int32 Type = Type{&C.ffi_type_sint32, "int32_t"}
	Int64 Type = Type{&C.ffi_type_sint64, "int64_t"}

	Float  Type = Type{&C.ffi_type_float, "float"}
	Double Type = Type{&C.ffi_type_double, "double"}

	Pointer Type = Type{&C.ffi_type_pointer, "pointer"}
)

func (t Type) String() string {
	if len(t.name) == 0 {
		return "struct"
	} else {
		return t.name
	}
}

type Interface struct {
	ffi_cif  C.ffi_cif
	ffi_ret  *C.ffi_type
	ffi_args **C.ffi_type

	ret  Type
	args []Type
}

func MustPrepare(ret Type, args ...Type) (cif Interface) {
	var err error

	if cif, err = Prepare(ret, args...); err != nil {
		panic(err)
	}

	return
}

func Prepare(ret Type, args ...Type) (cif Interface, err error) {
	cif.ffi_ret = ret.ffi_type
	cif.ret = ret
	cif.args = args
	argc := len(args)

	if argc != 0 {
		va := make([]*C.ffi_type, argc)

		for i, a := range args {
			va[i] = a.ffi_type
		}

		cif.ffi_args = &va[0]
	}

	if status := Status(C.ffi_prep_cif(&cif.ffi_cif, C.FFI_DEFAULT_ABI, C.uint(argc), cif.ffi_ret, cif.ffi_args)); status != OK {
		err = status
	}

	return
}

func (cif Interface) Call(fptr unsafe.Pointer, ret unsafe.Pointer, args ...unsafe.Pointer) (err error) {
	var va *unsafe.Pointer

	if len(args) != 0 {
		va = &args[0]
	}

	_, err = C.ffi_call(&cif.ffi_cif, C.function(fptr), ret, va)
	return
}

func Call(fptr unsafe.Pointer, ret interface{}, args ...interface{}) (err error) {
	vret := valueOfRet(ret)
	varg := valueOfArgs(args)

	rett := makeRetType(vret)
	retv := makeRetValue(vret)

	argt := makeArgTypes(varg)
	argv := makeArgValues(varg)

	err = MustPrepare(rett, argt...).Call(fptr, retv, argv...)

	setRetValue(vret, retv)
	return
}

func valueOfRet(ret interface{}) reflect.Value {
	v := reflect.ValueOf(ret)

	if ret != nil && v.Kind() != reflect.Ptr {
		panic(fmt.Sprintf("ffi: expected return value to be a pointer but got %T", ret))
	}

	return v
}

func valueOfArgs(args []interface{}) []reflect.Value {
	v := make([]reflect.Value, len(args))

	for i, a := range args {
		v[i] = reflect.ValueOf(a)
	}

	return v
}

func makeRetType(v reflect.Value) Type {
	if v.IsNil() {
		return Void
	}

	switch v.Elem().Kind() {
	case reflect.Int:
		return Int

	case reflect.Int8:
		return Int8

	case reflect.Int16:
		return Int16

	case reflect.Int32:
		return Int32

	case reflect.Int64:
		return Int64

	case reflect.Uint:
		return UInt

	case reflect.Uint8:
		return UInt8

	case reflect.Uint16:
		return UInt16

	case reflect.Uint32:
		return UInt32

	case reflect.Uint64:
		return UInt64

	case reflect.Float32:
		return Float

	case reflect.Float64:
		return Double

	case reflect.String, reflect.Ptr, reflect.UnsafePointer:
		return Pointer
	}

	unsupportedRetType(v)
	return Type{}
}

func makeRetValue(v reflect.Value) unsafe.Pointer {
	if v.IsNil() {
		return nil
	}

	switch v = v.Elem(); v.Kind() {
	case reflect.Int:
		x := C.int(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int8:
		x := C.int8_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int16:
		x := C.int16_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int32:
		x := C.int32_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int64:
		x := C.int64_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Uint8:
		x := C.uint8_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Uint16:
		x := C.uint16_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Uint32:
		x := C.uint32_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Uint64:
		x := C.uint64_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Float32:
		x := C.float(v.Float())
		return unsafe.Pointer(&x)

	case reflect.Float64:
		x := C.double(v.Float())
		return unsafe.Pointer(&x)

	case reflect.String, reflect.Ptr, reflect.UnsafePointer:
		x := unsafe.Pointer(nil)
		return unsafe.Pointer(&x)
	}

	unsupportedRetType(v)
	return nil
}

func makeArgTypes(v []reflect.Value) []Type {
	t := make([]Type, len(v))

	for i, a := range v {
		t[i] = makeArgType(a)
	}

	return t
}

func makeArgValues(v []reflect.Value) []unsafe.Pointer {
	p := make([]unsafe.Pointer, len(v))

	for i, a := range v {
		p[i] = makeArgValue(a)
	}

	return p
}

func makeArgType(v reflect.Value) Type {
	switch v.Kind() {
	case reflect.Int:
		return Int

	case reflect.Int8:
		return Int8

	case reflect.Int16:
		return Int16

	case reflect.Int32:
		return Int32

	case reflect.Int64:
		return Int64

	case reflect.Uint:
		return UInt

	case reflect.Uint8:
		return UInt8

	case reflect.Uint16:
		return UInt16

	case reflect.Uint32:
		return UInt32

	case reflect.Uint64:
		return UInt64

	case reflect.Float32:
		return Float

	case reflect.Float64:
		return Double

	case reflect.String, reflect.Slice, reflect.Ptr, reflect.UnsafePointer, reflect.Interface:
		return Pointer
	}

	unsupportedArgType(v)
	return Type{}
}

func makeArgValue(v reflect.Value) unsafe.Pointer {
	switch v.Kind() {
	case reflect.Int:
		x := C.int(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int8:
		x := C.int8_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int16:
		x := C.int16_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int32:
		x := C.int32_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Int64:
		x := C.int64_t(v.Int())
		return unsafe.Pointer(&x)

	case reflect.Uint8:
		x := C.uint8_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Uint16:
		x := C.uint16_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Uint32:
		x := C.uint32_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Uint64:
		x := C.uint64_t(v.Uint())
		return unsafe.Pointer(&x)

	case reflect.Float32:
		x := C.float(v.Float())
		return unsafe.Pointer(&x)

	case reflect.Float64:
		x := C.double(v.Float())
		return unsafe.Pointer(&x)

	case reflect.String:
		return unsafe.Pointer(C.CString(v.String()))

	case reflect.Slice, reflect.Ptr, reflect.UnsafePointer:
		return unsafe.Pointer(v.Pointer())

	case reflect.Interface:
		if v.IsNil() {
			return nil
		}
	}

	unsupportedArgType(v)
	return nil
}

func setRetValue(v reflect.Value, p unsafe.Pointer) {
	switch v = v.Elem(); v.Kind() {
	case reflect.Int:
		v.SetInt(int64(*((*C.int)(p))))

	case reflect.Int8:
		v.SetInt(int64(*(*C.int8_t)(p)))

	case reflect.Int16:
		v.SetInt(int64(*(*C.int16_t)(p)))

	case reflect.Int32:
		v.SetInt(int64(*(*C.int32_t)(p)))

	case reflect.Int64:
		v.SetInt(int64(*(*C.int64_t)(p)))

	case reflect.Uint:
		v.SetUint(uint64(*((*C.uint)(p))))

	case reflect.Uint8:
		v.SetUint(uint64(*(*C.uint8_t)(p)))

	case reflect.Uint16:
		v.SetUint(uint64(*(*C.uint16_t)(p)))

	case reflect.Uint32:
		v.SetUint(uint64(*(*C.uint32_t)(p)))

	case reflect.Uint64:
		v.SetUint(uint64(*(*C.uint64_t)(p)))

	case reflect.Float32:
		v.SetFloat(float64(*(*C.float)(p)))

	case reflect.Float64:
		v.SetFloat(float64(*(*C.double)(p)))

	case reflect.String:
		v.SetString(C.GoString(*(**C.char)(p)))

	case reflect.UnsafePointer, reflect.Ptr:
		v.SetPointer(*(*unsafe.Pointer)(p))
	}
}

func unsupportedArgType(v reflect.Value) {
	panic(fmt.Sprintf("ffi: unsupported argument type: %s", v.Type()))
}

func unsupportedRetType(v reflect.Value) {
	panic(fmt.Sprintf("ffi: unsupported return type: %s", v.Type()))
}

type Function interface {
	Pointer() uintptr
}

type function struct {
	Interface
	fptr unsafe.Pointer
	mptr unsafe.Pointer
	call reflect.Value
}

func (fn *function) Pointer() uintptr {
	return uintptr(fn.fptr)
}

func Closure(v interface{}) Function {
	switch f := v.(type) {
	case Function:
		return f
	}

	fv := reflect.ValueOf(v)
	ft := reflect.TypeOf(v)

	if ft.Kind() != reflect.Func {
		panic(fmt.Sprintf("ffi: closures can only be created from functions, got %s", ft))
	}

	if ft.IsVariadic() {
		panic(fmt.Sprintf("ffi: closures with a variable number of arguments are not supported"))
	}

	return makeClosure(fv, ft)
}

func makeClosure(fv reflect.Value, ft reflect.Type) *function {
	fn := &function{
		call: fv,
	}

	var rt Type
	var at []Type

	if n := ft.NumOut(); n != 0 {
		rt = makeRetType(reflect.New(ft.Out(0)))
	}

	if n := ft.NumIn(); n != 0 {
		at = make([]Type, n)

		for i := 0; i != n; i++ {
			at[i] = makeArgType(reflect.Zero(ft.In(i)))
		}
	}

	fn.Interface = MustPrepare(rt, at...)

	if err := constructClosure(fn); err != nil {
		panic(err)
	}

	runtime.SetFinalizer(fn, destroyClosure)
	return fn
}

//export GoClosureCallback
func GoClosureCallback(cif *C.ffi_cif, ret unsafe.Pointer, args *unsafe.Pointer, data unsafe.Pointer) {
	fn := (*function)(data)
	fv := fn.call
	ft := fv.Type()

	ac := ft.NumIn()
	av := make([]reflect.Value, ac)

	for i := 0; i != ac; i++ {
		av[i] = makeGoArg(*args, ft.In(i))
		args = nextUnsafePointer(args)
	}

	rv := fv.Call(av)
	rc := len(rv)

	if rc > 0 {
		setRetPointer(ret, rv[0])
	}

	if rc > 1 {
		// TODO: report errno
	}
}

func makeGoArg(p unsafe.Pointer, t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(*((*C.int)(p))))

	default:
		return reflect.ValueOf(nil)
	}
}

func setRetPointer(p unsafe.Pointer, v reflect.Value) {
	switch v.Kind() {
	case reflect.Int:
		*((*C.int)(p)) = C.int(v.Int())
	}
}

func nextUnsafePointer(p *unsafe.Pointer) *unsafe.Pointer {
	return (*unsafe.Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + uintptr(unsafe.Sizeof(*p))))
}
