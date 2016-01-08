package ffi

// #include <ffi.h>
//
// typedef void (*function)(void);
import "C"
import "unsafe"

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

type Type interface {
	ffi() *C.ffi_type
}

type ffi_type struct {
	ffi_type *C.ffi_type
}

func (t ffi_type) ffi() *C.ffi_type {
	return t.ffi_type
}

var (
	Void Type = ffi_type{&C.ffi_type_void}

	Uchar  Type = ffi_type{&C.ffi_type_schar}
	Ushort Type = ffi_type{&C.ffi_type_sshort}
	Uint   Type = ffi_type{&C.ffi_type_sint}
	Ulong  Type = ffi_type{&C.ffi_type_slong}

	Uint8  Type = ffi_type{&C.ffi_type_uint8}
	Uint16 Type = ffi_type{&C.ffi_type_uint16}
	Uint32 Type = ffi_type{&C.ffi_type_uint32}
	Uint64 Type = ffi_type{&C.ffi_type_uint64}

	Schar  Type = ffi_type{&C.ffi_type_schar}
	Sshort Type = ffi_type{&C.ffi_type_sshort}
	Sint   Type = ffi_type{&C.ffi_type_sint}
	SLong  Type = ffi_type{&C.ffi_type_slong}

	Sint8  Type = ffi_type{&C.ffi_type_uint8}
	Sint16 Type = ffi_type{&C.ffi_type_uint16}
	Sint32 Type = ffi_type{&C.ffi_type_uint32}
	Sint64 Type = ffi_type{&C.ffi_type_uint64}

	Float  Type = ffi_type{&C.ffi_type_float}
	Double Type = ffi_type{&C.ffi_type_double}

	Pointer Type = ffi_type{&C.ffi_type_pointer}
)

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
		ffi_ret: ret.ffi(),

		ret:  ret,
		args: args,
	}

	argc := len(args)

	if argc != 0 {
		va := make([]*C.ffi_type, argc)

		for i, a := range args {
			va[i] = a.ffi()
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

type Function func(...interface{}) (interface{}, error)

func Declare(fptr unsafe.Pointer, ret Type, args ...Type) Function {
	cif := Prepare(ret, args...)

	return func(argv ...interface{}) (res interface{}, err error) {

	}
}
