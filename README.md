go-ffi [![Build Status](https://travis-ci.org/achille-roussel/go-ffi.svg)](https://travis-ci.org/achille-roussel/go-ffi) [![Coverage Status](https://coveralls.io/repos/achille-roussel/go-ffi/badge.svg?branch=master&service=github)](https://coveralls.io/github/achille-roussel/go-ffi?branch=master)
======

Go bindings to libffi

Calling C Functions
-------------------

go-ffi is interesting when a program needs to call C functions in a way that is
not yet supported by cgo. For example if the function pointer was obtained
dynamically or if the function as a variadic signature.

Here's an example showing how to call printf from Go using the go-ffi package:
```go
package main

// #include <stdio.h>
//
// const void *printf__ = (void *) printf;
import "C"
import "github.com/achille-roussel/go-ffi"

func main() {
     ffi.Call(C.printf__, nil, "%s %s!\n", "Hello", "World")
}
```
```
Hello World!
```

Type Conversions
----------------

go-ffi automatically converts between C and Go types when using the high-level
interfaces.  
The following table exposes what conversions are supported:

| Go             | C               |
|----------------|-----------------|
| int            | int             |
| int8           | int8_t          |
| int16          | int16_t         |
| int32          | int32_t         |
| int64          | int64_t         |
| uint           | unsigned int    |
| uint8          | uint8_t         |
| uint16         | uint16_t        |
| uint32         | uint32_t        |
| uintptr        | size_t          |
| float32        | float           |
| float64        | double          |
| string         | char *          |
| unsafe.Pointer | void *          |

