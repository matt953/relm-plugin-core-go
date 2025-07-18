package ffi

/*
#include <stdlib.h>
#include <stdint.h>
#include <stdbool.h>
#include <string.h>

typedef struct {
    bool success;
    uint8_t* data;
    size_t data_len;
    char* error_msg;
} FFIResult;
*/
import "C"
import (
	"unsafe"
)

// FFIResult represents the result of an FFI operation
type FFIResult struct {
	Success  bool
	Data     []byte
	ErrorMsg string
}

// ToCFFIResult converts a Go FFIResult to a C FFIResult
func (r *FFIResult) ToCFFIResult() C.FFIResult {
	result := C.FFIResult{
		success: C.bool(r.Success),
	}

	if r.Success {
		if len(r.Data) > 0 {
			// Allocate C memory for data
			result.data = (*C.uint8_t)(C.malloc(C.size_t(len(r.Data))))
			result.data_len = C.size_t(len(r.Data))

			// Copy Go data to C memory
			C.memcpy(unsafe.Pointer(result.data), unsafe.Pointer(&r.Data[0]), result.data_len)
		} else {
			result.data = nil
			result.data_len = 0
		}
		result.error_msg = nil
	} else {
		result.data = nil
		result.data_len = 0
		if r.ErrorMsg != "" {
			result.error_msg = C.CString(r.ErrorMsg)
		} else {
			result.error_msg = nil
		}
	}

	return result
}

// NewSuccessResult creates a successful FFIResult with data
func NewSuccessResult(data []byte) *FFIResult {
	return &FFIResult{
		Success: true,
		Data:    data,
	}
}

// NewSuccessEmpty creates a successful FFIResult with no data
func NewSuccessEmpty() *FFIResult {
	return &FFIResult{
		Success: true,
		Data:    nil,
	}
}

// NewErrorResult creates an error FFIResult with a message
func NewErrorResult(errorMsg string) *FFIResult {
	return &FFIResult{
		Success:  false,
		ErrorMsg: errorMsg,
	}
}
