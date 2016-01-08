package ffi

// #include <ffi.h>
//
// typedef void (*function)(void);
import "C"
import (
	"fmt"
	"reflect"
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
}

var (
	Void Type = Type{&C.ffi_type_void}

	Uchar  Type = Type{&C.ffi_type_schar}
	Ushort Type = Type{&C.ffi_type_sshort}
	Uint   Type = Type{&C.ffi_type_sint}
	Ulong  Type = Type{&C.ffi_type_slong}

	Uint8  Type = Type{&C.ffi_type_uint8}
	Uint16 Type = Type{&C.ffi_type_uint16}
	Uint32 Type = Type{&C.ffi_type_uint32}
	Uint64 Type = Type{&C.ffi_type_uint64}

	Schar  Type = Type{&C.ffi_type_schar}
	Sshort Type = Type{&C.ffi_type_sshort}
	Sint   Type = Type{&C.ffi_type_sint}
	SLong  Type = Type{&C.ffi_type_slong}

	Sint8  Type = Type{&C.ffi_type_uint8}
	Sint16 Type = Type{&C.ffi_type_uint16}
	Sint32 Type = Type{&C.ffi_type_uint32}
	Sint64 Type = Type{&C.ffi_type_uint64}

	Float  Type = Type{&C.ffi_type_float}
	Double Type = Type{&C.ffi_type_double}

	Pointer Type = Type{&C.ffi_type_pointer}
)

func (t Type) String() string {
	switch t {
	case Void:
		return "void"

	case Uchar, Uint8:
		return "uint8"

	case Ushort, Uint16:
		return "uint16"

	case Uint, Uint32:
		return "uint32"

	case Ulong, Uint64:
		return "uint64"

	case Schar, Sint8:
		return "uint8"

	case Sshort, Sint16:
		return "uint16"

	case Sint, Sint32:
		return "uint32"

	case SLong, Sint64:
		return "uint64"

	case Float:
		return "float32"

	case Double:
		return "float64"

	case Pointer:
		return "pointer"

	default:
		return "struct"
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

	_, err = C.ffi_call(&cif.ffi_cif, (C.function)(fptr), ret, va)
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
	return makeType(v.Elem())
}

func makeRetValue(v reflect.Value) unsafe.Pointer {
	if v.IsNil() {
		return nil
	}
	return makeValue(v.Elem())
}

func makeArgTypes(v []reflect.Value) []Type {
	t := make([]Type, len(v))

	for i, a := range v {
		t[i] = makeType(a)
	}

	return t
}

func makeArgValues(v []reflect.Value) []unsafe.Pointer {
	p := make([]unsafe.Pointer, len(v))

	for i, a := range v {
		p[i] = makeValue(a)
	}

	return p
}

func makeType(v reflect.Value) Type {
	switch v.Kind() {
	case reflect.Int:
		return Sint
	}

	unsupportedType(v)
	return Type{}
}

func makeValue(v reflect.Value) unsafe.Pointer {
	switch v.Kind() {
	case reflect.Int:
		x := C.int(v.Int())
		return unsafe.Pointer(&x)
	}

	unsupportedType(v)
	return nil
}

func setRetValue(v reflect.Value, p unsafe.Pointer) {
	switch v = v.Elem(); v.Kind() {
	case reflect.Int:
		v.SetInt(int64(*((*C.int)(p))))
	}
}

func unsupportedType(v reflect.Value) {
	panic(fmt.Sprintf("ffi: unsupported type: %s", v.Type()))
}
