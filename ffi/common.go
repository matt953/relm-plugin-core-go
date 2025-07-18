package ffi

/*
#include <stdlib.h>
#include <string.h>
*/
import "C"
import (
	"unsafe"
)

// GoString converts a C string to a Go string
func GoString(cstr *C.char) string {
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

// CString converts a Go string to a C string
// Caller is responsible for freeing the returned C string
func CString(s string) *C.char {
	if s == "" {
		return nil
	}
	return C.CString(s)
}

// GoBytes converts C data to a Go byte slice
func GoBytes(data *C.uint8_t, length C.size_t) []byte {
	if data == nil || length == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}
