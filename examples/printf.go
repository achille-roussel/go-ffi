package main

// #include <stdio.h>
//
// const void *printf__ = (void *) printf;
import "C"
import "github.com/achille-roussel/go-ffi"

func main() {
	ffi.Call(C.printf__, nil, "%s %s!\n", "Hello", "World")
}
