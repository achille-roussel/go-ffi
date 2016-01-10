package ffi

// #cgo CFLAGS: -I/usr/include/ffi
// #cgo LDFLAGS: -lffi
//
// #include <ffi.h>
// #include <unistd.h>
// #include <sys/mman.h>
//
// static int ffi_closure_alloc__(void **mem) {
//   void *ptr;
//
//   ptr = mmap(NULL, sizeof(ffi_closure), PROT_READ | PROT_WRITE, MAP_ANON | MAP_PRIVATE, -1, 0);
//
//   if (MAP_FAILED == ptr) {
//     return -1;
//   }
//
//   *mem = ptr;
//   return 0;
// }
//
// static int ffi_closure_make_executable__(void *ptr) {
//   return mprotect(ptr, sizeof(ptr), PROT_READ | PROT_EXEC);
// }
//
// static void ffi_closure_free__(void *ptr) {
//   munmap(ptr, sizeof(ffi_closure));
// }
//
// extern void GoClosureCallback(ffi_cif *, void *, void **, void *);
//
// typedef  void (*closure)(ffi_cif*, void*, void**, void*);
//
import "C"
import "unsafe"

func constructClosure(fn *function) (err error) {
	var ptr unsafe.Pointer
	var rc C.int

	if rc, err = C.ffi_closure_alloc__(&ptr); rc != 0 {
		return
	}

	if status := Status(C.ffi_prep_closure((*C.ffi_closure)(ptr), &fn.Interface.ffi_cif, C.closure(C.GoClosureCallback), unsafe.Pointer(fn))); status != OK {
		C.ffi_closure_free__(ptr)
		err = status
		return
	}

	if rc, err = C.ffi_closure_make_executable__(ptr); rc != 0 {
		C.ffi_closure_free__(ptr)
		return
	}

	fn.fptr = ptr
	fn.mptr = ptr
	return nil
}

func destroyClosure(fn *function) {
	C.ffi_closure_free__(fn.mptr)
}
