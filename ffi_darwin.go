package ffi

// #cgo CFLAGS: -I/usr/include/ffi
// #cgo LDFLAGS: -lffi
//
// #include <ffi.h>
// #include <unistd.h>
// #include <sys/mman.h>
//
// static int ffi_closure_allocate__(void **addr) {
//   var *ptr;
//
//   ptr = mmap(NULL, sizeof(ffi_closure), PROT_READ | PROT_WRITE, MAP_ANON | MAP_PRIVATE, -1, 0);
//
//   if (MMAP_FAIL == ptr) {
//     return -1;
//   }
//
//   *addr = ptr;
//   return 0;
// }
//
// static int ffi_closure_make_executable__(void *addr) {
//    return mprotect(addr, sizeof(addr), PROT_READ | PROT_EXEC);
// }
//
// static void ffi_closure_free__(void *addr) {
//   munmap(addr, sizeof(ffi_closure));
// }
//
import "C"
import (
	"reflect"
	"runtime"
	"unsafe"
)

func ffi_closure_alloc(code *unsafe.Pointer) (mem unsafe.Pointer) {
	if C.ffi_closure_allocate__(&mem) != 0 {
		panic("ffi: closure allocation failed")
	}

	return
}

func makeClosure(fv reflect.Value, ft reflect.Type) *function {
	fn := &function{
		call: fv,
	}

	var cif Interface
	var rt Type
	var at []Type

	if n := ft.NumOut(); n != 0 {
		rt = makeRetType(reflect.Zero(ft.Out(0)))
	}

	if n := ft.NumIn(); n != 0 {
		at = make([]Type, n)

		for i := 0; i != n; i++ {
			at[i] = makeArgType(reflect.Zero(ft.In(i)))
		}
	}

	cif = MustPrepare(rt, at...)

	if fn.mptr = ffi_closure_alloc(unsafe.Sizeof(C.ffi_closure{}), &fn.fptr); fn.mptr == nil {
		panic("ffi: closure allocation failed")
	}

	if status := Status(C.ffi_prep_closure((*C.ffi_closure)(fn.mptr), &fn.ffi_cif, C.GoClosureCallback, unsafe.Pointer(fn), fn.fptr)); status != OK {
		C.ffi_closure_free(fn.mptr)
		panic("ffi: closure creation failed")
	}

	runtime.SetFinalizer(fn, (*function).destroy)
	return fn
}
