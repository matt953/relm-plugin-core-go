package general

import (
	"encoding/json"
	"fmt"

	"github.com/matt953/relm-types-go/types"
)

// Re-export Environment from the types library to avoid duplication
type Environment = types.Environment

// GeneralPlugin defines the interface for general-purpose plugins
// that handle various events and callbacks.
//
// General plugins provide event-driven functionality and can respond to various
// system events through callbacks. Unlike storage plugins which have specific
// data management functions, general plugins are more flexible and can handle
// any type of event or data.
type GeneralPlugin interface {
	// GetPluginName returns the human-readable name of the plugin
	GetPluginName() string

	// GetPluginVersion returns the version of the plugin
	GetPluginVersion() string

	// Initialize sets up the plugin (called once when loaded)
	// Returns true if initialization was successful
	Initialize() bool

	// Cleanup is called when the plugin is being unloaded
	// Returns true if cleanup was successful
	Cleanup() bool

	// OnEnvironmentCreated is called when a new environment is created
	//
	// This method provides the abstracted interface - developers just need to implement
	// this method and receive the parsed Environment directly.
	//
	// Parameters:
	//   environment - The parsed environment data
	//
	// Returns:
	//   true if the callback was processed successfully
	//   false if there was a processing error
	OnEnvironmentCreated(environment *Environment) bool
}

// CallPluginCallback is a helper function to call a plugin callback with serialized data
//
// This is a convenience function for plugin implementations that handles
// JSON serialization and error handling. It can be used for any callback
// that takes serializable data.
//
// Parameters:
//   callbackName - Name of the callback for error reporting
//   data - The data to serialize and pass to the callback
//   callbackFn - The actual callback function to call
//
// Returns:
//   true if the callback succeeded
//   false if the callback failed
//
// Example:
//
//	type MyData struct {
//	    ID   string `json:"id"`
//	    Name string `json:"name"`
//	}
//
//	func (p *MyPlugin) callMyCallback(data MyData) bool {
//	    return CallPluginCallback("my_custom_callback", data, func(jsonStr string) bool {
//	        fmt.Printf("Custom callback with data: %s\n", jsonStr)
//	        return true
//	    })
//	}
func CallPluginCallback[T any](callbackName string, data T, callbackFn func(string) bool) bool {
	// Serialize the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("Failed to serialize data for %s: %v\n", callbackName, err)
		return false
	}

	// Call the callback function with the JSON data
	return callbackFn(string(jsonData))
}

// MustCallPluginCallback is like CallPluginCallback but panics on serialization errors
//
// This version should only be used when you are certain that the data can be serialized
// successfully, such as with well-defined structs.
//
// Parameters:
//   callbackName - Name of the callback for error reporting
//   data - The data to serialize and pass to the callback
//   callbackFn - The actual callback function to call
//
// Returns:
//   true if the callback succeeded
//   false if the callback failed
//
// Panics:
//   If JSON serialization fails
func MustCallPluginCallback[T any](callbackName string, data T, callbackFn func(string) bool) bool {
	// Serialize the data to JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		panic(fmt.Sprintf("Failed to serialize data for %s: %v", callbackName, err))
	}

	// Call the callback function with the JSON data
	return callbackFn(string(jsonData))
}

// ParseCallbackData is a helper function to parse JSON callback data
//
// This function helps plugin implementations parse incoming JSON data
// into strongly-typed Go structs.
//
// Parameters:
//   jsonData - The JSON string to parse
//   target - Pointer to the target struct to unmarshal into
//
// Returns:
//   error if parsing fails, nil on success
//
// Example:
//
//	func (p *MyPlugin) OnEnvironmentCreated(environmentJSON string) bool {
//	    var env types.Environment
//	    if err := ParseCallbackData(environmentJSON, &env); err != nil {
//	        fmt.Printf("Error parsing environment data: %v\n", err)
//	        return false
//	    }
//	    
//	    fmt.Printf("Environment created: %s (%s)\n", env.Name, env.ID)
//	    return true
//	}
func ParseCallbackData[T any](jsonData string, target *T) error {
	return json.Unmarshal([]byte(jsonData), target)
}