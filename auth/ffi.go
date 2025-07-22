package auth

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
	"encoding/json"
	"runtime"
	"runtime/debug"
	"unsafe"

	"github.com/matt953/relm-plugin-core-go/config"
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
		success:   true,
		error_msg: nil,
	}

	if len(data) > 0 {
		result.data = (*C.uint8_t)(C.malloc(C.size_t(len(data))))
		result.data_len = C.size_t(len(data))
		C.memcpy(unsafe.Pointer(result.data), unsafe.Pointer(&data[0]), C.size_t(len(data)))
	} else {
		result.data = nil
		result.data_len = 0
	}

	return result
}

func newSuccessJSONResult(v interface{}) C.FFIResult {
	data, err := json.Marshal(v)
	if err != nil {
		return newErrorResult("Failed to serialize response: " + err.Error())
	}
	return newSuccessResult(data)
}

func newErrorResult(msg string) C.FFIResult {
	return C.FFIResult{
		success:   false,
		data:      nil,
		data_len:  0,
		error_msg: cString(msg),
	}
}

func newSuccessEmptyResult() C.FFIResult {
	return C.FFIResult{
		success:   true,
		data:      nil,
		data_len:  0,
		error_msg: nil,
	}
}

// Exported C functions for auth plugins

//export initialize_with_config
func initialize_with_config(configJson *C.char) C.bool {
	configStr := goString(configJson)

	// Set the global config
	if err := config.SetConfigFromJSON(configStr); err != nil {
		println("initialize_with_config: failed to set global config:", err.Error())
		return C.bool(false)
	}

	// Plugin initialization - verify we have a registered plugin
	plugin := GetRegisteredPlugin()
	return C.bool(plugin != nil)
}

//export check_user_access
func check_user_access(userID, resource, action *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if userID == nil || resource == nil || action == nil {
		return newErrorResult("userID, resource, and action cannot be null")
	}

	userIDStr := goString(userID)
	resourceStr := goString(resource)
	actionStr := goString(action)

	allowed, err := plugin.CheckUserAccess(userIDStr, resourceStr, actionStr)
	if err != nil {
		return newErrorResult("Failed to check user access: " + err.Error())
	}

	return newSuccessJSONResult(allowed)
}

//export get_user_permissions
func get_user_permissions(userID *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if userID == nil {
		return newErrorResult("userID cannot be null")
	}

	userIDStr := goString(userID)
	permissions, err := plugin.GetUserPermissions(userIDStr)
	if err != nil {
		return newErrorResult("Failed to get user permissions: " + err.Error())
	}

	return newSuccessJSONResult(permissions)
}

//export get_user_details
func get_user_details(userID *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if userID == nil {
		return newErrorResult("userID cannot be null")
	}

	userIDStr := goString(userID)
	details, err := plugin.GetUserDetails(userIDStr)
	if err != nil {
		return newErrorResult("Failed to get user details: " + err.Error())
	}

	return newSuccessJSONResult(details)
}

//export get_plugin_info
func get_plugin_info() C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	info, err := plugin.GetPluginInfo()
	if err != nil {
		return newErrorResult("Failed to get plugin info: " + err.Error())
	}

	return newSuccessJSONResult(info)
}

//export provider_name
func provider_name() *C.char {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return cString("unknown")
	}

	return cString(plugin.ProviderName())
}

//export health_check
func health_check() bool {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return false
	}

	return plugin.HealthCheck()
}

//export validate_user
func validate_user(userID *C.char) bool {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return false
	}

	if userID == nil {
		return false
	}

	userIDStr := goString(userID)
	return plugin.ValidateUser(userIDStr)
}

//export get_user_groups
func get_user_groups(userID *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if userID == nil {
		return newErrorResult("userID cannot be null")
	}

	userIDStr := goString(userID)
	groups, err := plugin.GetUserGroups(userIDStr)
	if err != nil {
		return newErrorResult("Failed to get user groups: " + err.Error())
	}

	return newSuccessJSONResult(groups)
}

//export search_users
func search_users(query *C.char, limit C.size_t) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if query == nil {
		return newErrorResult("query cannot be null")
	}

	queryStr := goString(query)
	limitInt := int(limit)

	users, err := plugin.SearchUsers(queryStr, limitInt)
	if err != nil {
		return newErrorResult("Failed to search users: " + err.Error())
	}

	return newSuccessJSONResult(users)
}

//export cleanup_plugin
func cleanup_plugin() bool {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return false
	}

	err := plugin.Cleanup()
	return err == nil
}

// Force GC to run periodically
func init() {
	go func() {
		runtime.GC()
		debug.FreeOSMemory()
	}()
}

func main() {
	// Required for CGO plugins
}
