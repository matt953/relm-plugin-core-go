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
	"github.com/matt953/relm-types-go/types"
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

//export create_user
func create_user(request *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if request == nil {
		return newErrorResult("request cannot be null")
	}

	requestStr := goString(request)

	var createRequest types.CreateUserRequest
	if err := json.Unmarshal([]byte(requestStr), &createRequest); err != nil {
		return newErrorResult("Failed to parse request JSON: " + err.Error())
	}

	userDetails, err := plugin.CreateUser(createRequest)
	if err != nil {
		return newErrorResult("Failed to create user: " + err.Error())
	}

	return newSuccessJSONResult(userDetails)
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

//export get_user_details_by_email
func get_user_details_by_email(email *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if email == nil {
		return newErrorResult("email cannot be null")
	}

	emailStr := goString(email)
	details, err := plugin.GetUserDetailsByEmail(emailStr)
	if err != nil {
		return newErrorResult("Failed to get user details by email: " + err.Error())
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

//export delete_user
func delete_user(userID *C.char) C.FFIResult {
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
	err := plugin.DeleteUser(userIDStr)
	if err != nil {
		return newErrorResult("Failed to delete user: " + err.Error())
	}

	return newSuccessEmptyResult()
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

// OAuth Client Management FFI exports

//export create_oauth_client
func create_oauth_client(request *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if request == nil {
		return newErrorResult("request cannot be null")
	}

	requestStr := goString(request)

	var createRequest types.CreateOAuthClientRequest
	if err := json.Unmarshal([]byte(requestStr), &createRequest); err != nil {
		return newErrorResult("Failed to parse request JSON: " + err.Error())
	}

	client, err := plugin.CreateOAuthClient(createRequest)
	if err != nil {
		return newErrorResult("Failed to create OAuth client: " + err.Error())
	}

	return newSuccessJSONResult(client)
}

//export get_oauth_client
func get_oauth_client(clientID *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if clientID == nil {
		return newErrorResult("clientID cannot be null")
	}

	clientIDStr := goString(clientID)
	client, err := plugin.GetOAuthClient(clientIDStr)
	if err != nil {
		return newErrorResult("Failed to get OAuth client: " + err.Error())
	}

	return newSuccessJSONResult(client)
}

//export update_oauth_client
func update_oauth_client(clientID, request *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if clientID == nil || request == nil {
		return newErrorResult("clientID and request cannot be null")
	}

	clientIDStr := goString(clientID)
	requestStr := goString(request)

	var updateRequest types.UpdateOAuthClientRequest
	if err := json.Unmarshal([]byte(requestStr), &updateRequest); err != nil {
		return newErrorResult("Failed to parse request JSON: " + err.Error())
	}

	client, err := plugin.UpdateOAuthClient(clientIDStr, updateRequest)
	if err != nil {
		return newErrorResult("Failed to update OAuth client: " + err.Error())
	}

	return newSuccessJSONResult(client)
}

//export delete_oauth_client
func delete_oauth_client(clientID *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if clientID == nil {
		return newErrorResult("clientID cannot be null")
	}

	clientIDStr := goString(clientID)
	err := plugin.DeleteOAuthClient(clientIDStr)
	if err != nil {
		return newErrorResult("Failed to delete OAuth client: " + err.Error())
	}

	return newSuccessEmptyResult()
}

//export list_oauth_clients
func list_oauth_clients(limit, offset C.int) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	var limitPtr, offsetPtr *int
	if limit >= 0 {
		limitInt := int(limit)
		limitPtr = &limitInt
	}
	if offset >= 0 {
		offsetInt := int(offset)
		offsetPtr = &offsetInt
	}

	clients, err := plugin.ListOAuthClients(limitPtr, offsetPtr)
	if err != nil {
		return newErrorResult("Failed to list OAuth clients: " + err.Error())
	}

	return newSuccessJSONResult(clients)
}

//export list_user_authorized_clients
func list_user_authorized_clients(userID *C.char) C.FFIResult {
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
	clients, err := plugin.ListUserAuthorizedClients(userIDStr)
	if err != nil {
		return newErrorResult("Failed to list user authorized clients: " + err.Error())
	}

	return newSuccessJSONResult(clients)
}

//export revoke_user_client_authorization
func revoke_user_client_authorization(userID, clientID *C.char) C.FFIResult {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
		}
	}()

	plugin := GetRegisteredPlugin()
	if plugin == nil {
		return newErrorResult("No plugin registered")
	}

	if userID == nil || clientID == nil {
		return newErrorResult("userID and clientID cannot be null")
	}

	userIDStr := goString(userID)
	clientIDStr := goString(clientID)
	err := plugin.RevokeUserClientAuthorization(userIDStr, clientIDStr)
	if err != nil {
		return newErrorResult("Failed to revoke user client authorization: " + err.Error())
	}

	return newSuccessEmptyResult()
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
