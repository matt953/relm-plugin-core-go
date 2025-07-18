package storage

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
	"runtime"
	"runtime/debug"
	"unsafe"
)

// Helper functions for C interop
func goString(cstr *C.char) string {
	if cstr == nil {
		return ""
	}
	return C.GoString(cstr)
}

func cString(s string) *C.char {
	if s == "" {
		return nil
	}
	return C.CString(s)
}

func goBytes(data *C.uint8_t, length C.size_t) []byte {
	if data == nil || length == 0 {
		return nil
	}
	return C.GoBytes(unsafe.Pointer(data), C.int(length))
}

// FFI result helpers
func newSuccessResult(data []byte) C.FFIResult {
	result := C.FFIResult{
		success: C.bool(true),
	}

	if len(data) > 0 {
		// Allocate C memory for data
		result.data = (*C.uint8_t)(C.malloc(C.size_t(len(data))))
		result.data_len = C.size_t(len(data))

		// Copy Go data to C memory
		C.memcpy(unsafe.Pointer(result.data), unsafe.Pointer(&data[0]), result.data_len)
	} else {
		result.data = nil
		result.data_len = 0
	}
	result.error_msg = nil

	return result
}

func newSuccessEmpty() C.FFIResult {
	return C.FFIResult{
		success:   C.bool(true),
		data:      nil,
		data_len:  0,
		error_msg: nil,
	}
}

func newErrorResult(errorMsg string) C.FFIResult {
	var error_msg *C.char
	if errorMsg != "" {
		error_msg = C.CString(errorMsg)
	} else {
		error_msg = nil
	}

	return C.FFIResult{
		success:   C.bool(false),
		data:      nil,
		data_len:  0,
		error_msg: error_msg,
	}
}

// FFI export functions that will be called from Rust

//export store_file_with_content_type
func store_file_with_content_type(path *C.char, data *C.uint8_t, length C.size_t, contentType *C.char) C.FFIResult {
	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	goPath := goString(path)
	goData := goBytes(data, length)

	var goContentType *string
	if contentType != nil {
		ct := goString(contentType)
		goContentType = &ct
	}

	err := plugin.StoreFile(goPath, goData, goContentType)
	if err != nil {
		return newErrorResult(err.Error())
	}

	return newSuccessEmpty()
}

//export store_file
func store_file(path *C.char, data *C.uint8_t, length C.size_t) C.FFIResult {
	return store_file_with_content_type(path, data, length, nil)
}

//export retrieve_file
func retrieve_file(path *C.char) C.FFIResult {
	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	goPath := goString(path)

	data, err := plugin.RetrieveFile(goPath)
	if err != nil {
		return newErrorResult(err.Error())
	}

	return newSuccessResult(data)
}

//export delete_file
func delete_file(path *C.char) C.FFIResult {
	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	goPath := goString(path)

	err := plugin.DeleteFile(goPath)
	if err != nil {
		return newErrorResult(err.Error())
	}

	return newSuccessEmpty()
}

//export file_exists
func file_exists(path *C.char) C.bool {
	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return C.bool(false)
	}

	goPath := goString(path)
	exists := plugin.FileExists(goPath)
	return C.bool(exists)
}

//export generate_file_url
func generate_file_url(path *C.char, baseURL *C.char) *C.char {
	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return nil
	}

	goPath := goString(path)
	goBaseURL := goString(baseURL)

	url := plugin.GenerateURL(goPath, goBaseURL)
	if url == nil {
		return nil
	}

	return cString(*url)
}

//export provider_name
func provider_name() *C.char {
	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return cString("Unknown Go Plugin")
	}

	name := plugin.ProviderName()
	return cString(name)
}

//export init_plugin
func init_plugin() C.bool {
	// Plugin initialization - just verify we have a registered plugin
	plugin := GetRegisteredPlugin()
	return C.bool(plugin != nil)
}

//export cleanup_plugin
func cleanup_plugin() C.bool {
	// Add debug logging to see if this is being called
	println("cleanup_plugin called")

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		println("cleanup_plugin: no plugin registered")
		return C.bool(false)
	}

	err := plugin.Cleanup()
	if err != nil {
		println("cleanup_plugin: plugin cleanup failed:", err.Error())
		return C.bool(false)
	}

	// Force garbage collection to clean up any remaining resources
	runtime.GC()
	debug.FreeOSMemory()

	println("cleanup_plugin: completed successfully")

	return C.bool(true)
}

// Required for CGO exports
func main() {}
