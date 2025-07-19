package general

/*
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"fmt"
	"sync"
	"unsafe"
)

var (
	pluginInstance GeneralPlugin
	pluginMutex    sync.RWMutex
	isInitialized  bool
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
//   plugin - The GeneralPlugin implementation to export
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

//export free_string
func free_string(ptr *C.char) {
	if ptr != nil {
		C.free(unsafe.Pointer(ptr))
	}
}

// Required for CGO
func main() {}