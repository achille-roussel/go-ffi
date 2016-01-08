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

type Interface interface {
	Call(fptr unsafe.Pointer, ret unsafe.Pointer, args ...unsafe.Pointer) error
}

type ffi_cif struct {
	ffi_cif  C.ffi_cif
	ffi_ret  *C.ffi_type
	ffi_args **C.ffi_type

	ret  Type
	args []Type
}

func Prepare(ret Type, args ...Type) (Interface, error) {
	cif := ffi_cif{
		ffi_ret: ret.ffi_type,

		ret:  ret,
		args: args,
	}

	argc := len(args)

	if argc != 0 {
		va := make([]*C.ffi_type, argc)

		for i, a := range args {
			va[i] = a.ffi_type
		}

		cif.ffi_args = &va[0]
	}

	if status := Status(C.ffi_prep_cif(&cif.ffi_cif, C.FFI_DEFAULT_ABI, C.uint(argc), cif.ffi_ret, cif.ffi_args)); status != OK {
		return nil, status
	}

	return cif, nil
}

func (cif ffi_cif) Call(fptr unsafe.Pointer, ret unsafe.Pointer, args ...unsafe.Pointer) (err error) {
	var va *unsafe.Pointer

	if len(args) != 0 {
		va = &args[0]
	}

	_, err = C.ffi_call(&cif.ffi_cif, (C.function)(fptr), ret, va)
	return
}

func Declare(fptr unsafe.Pointer, ret Type, args ...Type) interface{} {
	var cif Interface
	var err error

	if cif, err = Prepare(ret, args...); err != nil {
		panic(err)
	}

	size := callSizeOf(ret, args)
	sig := signatureOf(ret, args)

	fun := func(argv []reflect.Value) (res []reflect.Value) {
		alloc := makeAllocator(size)

		vr := alloc.allocate(sizeOf(ret))
		va := makeArgs(argv, args, &alloc)
		defer freeArgs(argv, va)

		// err := cif.Call(fptr, vr, va...)
		var err error
		_ = cif

		switch sig.NumOut() {
		case 1:
			res = []reflect.Value{reflect.ValueOf(err)}
		default:
			res = []reflect.Value{makeRet(sig.Out(0), vr), reflect.ValueOf(err)}
		}

		return
	}

	return reflect.MakeFunc(sig, fun).Interface()
}

func signatureOf(ret Type, args []Type) reflect.Type {
	return reflect.FuncOf(argsTypeOf(args), retTypeOf(ret), false)
}

func argsTypeOf(args []Type) []reflect.Type {
	var ts = make([]reflect.Type, len(args))

	for i, a := range args {
		ts[i] = typeOf(a)
	}

	return ts
}

func retTypeOf(ret Type) []reflect.Type {
	var err error
	var ts = make([]reflect.Type, 0, 2)

	if t := typeOf(ret); t != nil {
		ts = append(ts, t)
	}

	return append(ts, reflect.TypeOf(&err).Elem())
}

func typeOf(t Type) reflect.Type {
	switch t {
	case Void:
		return nil

	case Sint:
		return reflect.TypeOf(int(0))

		// TODO: add more types
	default:
		panic(fmt.Sprintf("ffi: unsupported type: %#v", t))
	}
}

func callSizeOf(ret Type, args []Type) int {
	size := sizeOf(ret)

	for _, a := range args {
		size += sizeOf(a)
	}

	return size
}

func sizeOf(t Type) int {
	return align(int(t.ffi_type.size), int(t.ffi_type.alignment))
}

func align(size int, alignment int) int {
	n := size / alignment

	if (size % alignment) != 0 {
		n++
	}

	return n
}

func makeArgs(argv []reflect.Value, args []Type, alloc *allocator) []unsafe.Pointer {
	va := make([]unsafe.Pointer, len(args))

	for i, a := range args {
		va[i] = makeArg(argv[i], alloc.allocate(sizeOf(a)))
	}

	return va
}

func makeArg(v reflect.Value, p unsafe.Pointer) unsafe.Pointer {
	switch v.Kind() {
	case reflect.Int:
		*((*C.int)(p)) = C.int(v.Int())

		// TODO: add more types
	default:
		panic("ffi: an unreachble portion of code is executed: mismatched argument types between Go and C")
	}

	return p
}

func freeArgs(argv []reflect.Value, va []unsafe.Pointer) {
	for i, v := range argv {
		freeArg(v, va[i])
	}
}

func freeArg(v reflect.Value, p unsafe.Pointer) {
	// TODO: free if necessary
}

func makeRet(t reflect.Type, p unsafe.Pointer) (v reflect.Value) {
	v = reflect.New(t).Elem()

	switch t.Kind() {
	case reflect.Int:
		v.SetInt(int64(*((*C.int)(p))))

		// TODO: add more types
	default:
		panic("ffi: an unreachble portion of code is executed: mismatched argument types between Go and C")
	}

	return
}

type allocator struct {
	bytes []byte
	addr  uintptr
}

func makeAllocator(n int) (a allocator) {
	if n != 0 {
		a.bytes = make([]byte, n)
		a.addr = uintptr(unsafe.Pointer(&a.bytes[0]))
	}
	return
}

func (a *allocator) allocate(n int) unsafe.Pointer {
	ptr := unsafe.Pointer(a.addr)
	a.addr += uintptr(n)
	return ptr
}
