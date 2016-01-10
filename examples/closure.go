package main

import (
	"fmt"
	"strconv"
	"unsafe"

	"github.com/achille-roussel/go-ffi"
)

func main() {
	itoa := ffi.Closure(strconv.Itoa)
	fptr := itoa.Pointer()
	repr := ""

	ffi.Call(unsafe.Pointer(fptr), &repr, 42)

	fmt.Println(repr)
}
