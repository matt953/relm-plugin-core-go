package general

/*
#include <stdlib.h>
#include <stdbool.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"

	"github.com/matt953/relm-plugin-core-go/config"
)

var (
	pluginInstance GeneralPlugin
	pluginMutex    sync.RWMutex
	isInitialized  bool
	pluginName     string
)

// ExportPlugin registers a GeneralPlugin implementation for FFI export
//
// This function should be called from your plugin's main function to register
// the plugin implementation and make it available to the FFI system.
//
// Example:
//
//	func main() {
//	    plugin := &MyGeneralPlugin{}
//	    ExportPlugin(plugin)
//	}
//
// Parameters:
//
//	plugin - The GeneralPlugin implementation to export
func ExportPlugin(plugin GeneralPlugin) {
	pluginMutex.Lock()
	defer pluginMutex.Unlock()

	pluginInstance = plugin
	fmt.Printf("General plugin exported: %s v%s\n",
		plugin.GetPluginName(), plugin.GetPluginVersion())
}

// getPlugin safely returns the current plugin instance
func getPlugin() GeneralPlugin {
	pluginMutex.RLock()
	defer pluginMutex.RUnlock()
	return pluginInstance
}

//export get_plugin_name
func get_plugin_name() *C.char {
	plugin := getPlugin()
	if plugin == nil {
		return C.CString("Unknown Plugin")
	}
	return C.CString(plugin.GetPluginName())
}

//export get_plugin_version
func get_plugin_version() *C.char {
	plugin := getPlugin()
	if plugin == nil {
		return C.CString("0.0.0")
	}
	return C.CString(plugin.GetPluginVersion())
}

//export initialize_with_config
func initialize_with_config(configJson *C.char) C.bool {
	configStr := C.GoString(configJson)

	plugin := getPlugin()
	if plugin == nil {
		return C.bool(false)
	}

	pluginMutex.Lock()
	defer pluginMutex.Unlock()

	// Set the global config
	if err := config.SetConfigFromJSON(configStr); err != nil {
		fmt.Printf("initialize_with_config: failed to set global config: %v\n", err)
		return C.bool(false)
	}

	if isInitialized {
		return C.bool(true) // Already initialized but config is now set
	}

	success := plugin.Initialize()
	if success {
		isInitialized = true
		return C.bool(true)
	}
	return C.bool(false)
}

//export init_plugin
func init_plugin() C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	pluginMutex.Lock()
	defer pluginMutex.Unlock()

	if isInitialized {
		return 1 // Already initialized
	}


	success := plugin.Initialize()
	if success {
		isInitialized = true
		return 1
	}
	return 0
}

//export cleanup_plugin
func cleanup_plugin() C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	pluginMutex.Lock()
	defer pluginMutex.Unlock()

	if !isInitialized {
		return 1 // Already cleaned up
	}

	success := plugin.Cleanup()
	if success {
		isInitialized = false
		return 1
	}
	return 0
}

//export on_environment_created
func on_environment_created(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	// Parse the JSON automatically for the plugin developer
	var environment Environment
	if err := json.Unmarshal([]byte(jsonStr), &environment); err != nil {
		return 0
	}

	success := plugin.OnEnvironmentCreated(&environment)
	if success {
		return 1
	}
	return 0
}

//export on_environment_updated
func on_environment_updated(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var environment Environment
	if err := json.Unmarshal([]byte(jsonStr), &environment); err != nil {
		return 0
	}

	success := plugin.OnEnvironmentUpdated(&environment)
	if success {
		return 1
	}
	return 0
}

//export on_environment_deleted
func on_environment_deleted(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var environment Environment
	if err := json.Unmarshal([]byte(jsonStr), &environment); err != nil {
		return 0
	}

	success := plugin.OnEnvironmentDeleted(&environment)
	if success {
		return 1
	}
	return 0
}

//export on_organization_created
func on_organization_created(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var organization Organization
	if err := json.Unmarshal([]byte(jsonStr), &organization); err != nil {
		return 0
	}

	success := plugin.OnOrganizationCreated(&organization)
	if success {
		return 1
	}
	return 0
}

//export on_organization_updated
func on_organization_updated(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var organization Organization
	if err := json.Unmarshal([]byte(jsonStr), &organization); err != nil {
		return 0
	}

	success := plugin.OnOrganizationUpdated(&organization)
	if success {
		return 1
	}
	return 0
}

//export on_organization_deleted
func on_organization_deleted(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var organization Organization
	if err := json.Unmarshal([]byte(jsonStr), &organization); err != nil {
		return 0
	}

	success := plugin.OnOrganizationDeleted(&organization)
	if success {
		return 1
	}
	return 0
}

//export on_user_created
func on_user_created(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var user User
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		return 0
	}

	success := plugin.OnUserCreated(&user)
	if success {
		return 1
	}
	return 0
}

//export on_user_updated
func on_user_updated(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var user User
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		return 0
	}

	success := plugin.OnUserUpdated(&user)
	if success {
		return 1
	}
	return 0
}

//export on_user_deleted
func on_user_deleted(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var user User
	if err := json.Unmarshal([]byte(jsonStr), &user); err != nil {
		return 0
	}

	success := plugin.OnUserDeleted(&user)
	if success {
		return 1
	}
	return 0
}

//export on_user_assigned_to_organization
func on_user_assigned_to_organization(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var params struct {
		UserID string `json:"user_id"`
		OrgID  string `json:"org_id"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return 0
	}

	success := plugin.OnUserAssignedToOrganization(params.UserID, params.OrgID)
	if success {
		return 1
	}
	return 0
}

//export on_user_removed_from_organization
func on_user_removed_from_organization(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var params struct {
		UserID string `json:"user_id"`
		OrgID  string `json:"org_id"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return 0
	}

	success := plugin.OnUserRemovedFromOrganization(params.UserID, params.OrgID)
	if success {
		return 1
	}
	return 0
}

//export on_user_assigned_to_environment
func on_user_assigned_to_environment(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var params struct {
		UserID string `json:"user_id"`
		EnvID  string `json:"env_id"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return 0
	}

	success := plugin.OnUserAssignedToEnvironment(params.UserID, params.EnvID)
	if success {
		return 1
	}
	return 0
}

//export on_user_removed_from_environment
func on_user_removed_from_environment(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var params struct {
		UserID string `json:"user_id"`
		EnvID  string `json:"env_id"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return 0
	}

	success := plugin.OnUserRemovedFromEnvironment(params.UserID, params.EnvID)
	if success {
		return 1
	}
	return 0
}

//export on_user_avatar_uploaded
func on_user_avatar_uploaded(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var params struct {
		UserID string `json:"user_id"`
		FileID string `json:"file_id"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &params); err != nil {
		return 0
	}

	success := plugin.OnUserAvatarUploaded(params.UserID, params.FileID)
	if success {
		return 1
	}
	return 0
}

// OAuth Client Callback FFI Functions

//export on_oauth_client_created
func on_oauth_client_created(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var event OAuthClientCreatedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		return 0
	}

	success := plugin.OnOAuthClientCreated(&event)
	if success {
		return 1
	}
	return 0
}

//export on_oauth_client_updated
func on_oauth_client_updated(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var event OAuthClientUpdatedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		return 0
	}

	success := plugin.OnOAuthClientUpdated(&event)
	if success {
		return 1
	}
	return 0
}

//export on_oauth_client_deleted
func on_oauth_client_deleted(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var event OAuthClientDeletedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		return 0
	}

	success := plugin.OnOAuthClientDeleted(&event)
	if success {
		return 1
	}
	return 0
}

//export on_user_client_authorization_revoked
func on_user_client_authorization_revoked(jsonPtr *C.char) C.int {
	plugin := getPlugin()
	if plugin == nil {
		return 0
	}

	if jsonPtr == nil {
		return 0
	}

	jsonStr := C.GoString(jsonPtr)

	var event UserClientAuthorizationRevokedEvent
	if err := json.Unmarshal([]byte(jsonStr), &event); err != nil {
		return 0
	}

	success := plugin.OnUserClientAuthorizationRevoked(&event)
	if success {
		return 1
	}
	return 0
}

//export free_string
func free_string(ptr *C.char) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

// Required for CGO
func main() {}
