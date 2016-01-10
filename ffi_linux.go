package ffi

// #cgo LDFLAGS: -lffi
//
// #include <ffi.h>
//
// extern void GoClosureCallback(ffi_cif *, void *, void **, void *);
//
// typedef void (*closure)(ffi_cif *, void *, void **, void *);
import "C"
import "unsafe"

func constructClosure(fn *function) (err error) {
	var closure C.ffi_closure
	var fptr unsafe.Pointer
	var mptr unsafe.Pointer

	if mptr, err = C.ffi_closure_alloc(C.size_t(unsafe.Sizeof(closure)), &fptr); mptr == nil {
		return
	}

	if status := Status(C.ffi_prep_closure_loc((*C.ffi_closure)(mptr), &fn.Interface.ffi_cif, C.closure(C.GoClosureCallback), unsafe.Pointer(fn), fptr)); status != OK {
		C.ffi_closure_free(mptr)
		err = status
		return
	}

	fn.fptr = fptr
	fn.mptr = mptr
	return nil
}

func destroyClosure(fn *function) {
	C.ffi_closure_free(fn.mptr)
}
