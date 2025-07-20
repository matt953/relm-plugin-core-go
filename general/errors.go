package general

import "fmt"

// PluginError represents errors that can occur in general plugins
type PluginError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error implements the error interface
func (e *PluginError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Type, e.Message, e.Details)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Common error types for general plugins
var (
	// ErrInvalidInput indicates that the input data provided to a function is invalid
	ErrInvalidInput = func(msg string, details ...string) *PluginError {
		var detail string
		if len(details) > 0 {
			detail = details[0]
		}
		return &PluginError{
			Type:    "InvalidInput",
			Message: msg,
			Details: detail,
		}
	}

	// ErrCallbackFailed indicates that a callback function failed to execute
	ErrCallbackFailed = func(callbackName string, details ...string) *PluginError {
		var detail string
		if len(details) > 0 {
			detail = details[0]
		}
		return &PluginError{
			Type:    "CallbackFailed",
			Message: fmt.Sprintf("Callback '%s' failed", callbackName),
			Details: detail,
		}
	}

	// ErrNotInitialized indicates that a plugin operation was attempted before initialization
	ErrNotInitialized = func(pluginName string) *PluginError {
		return &PluginError{
			Type:    "NotInitialized",
			Message: fmt.Sprintf("Plugin '%s' is not initialized", pluginName),
		}
	}

	// ErrConfigurationError indicates a problem with plugin configuration
	ErrConfigurationError = func(msg string, details ...string) *PluginError {
		var detail string
		if len(details) > 0 {
			detail = details[0]
		}
		return &PluginError{
			Type:    "ConfigurationError",
			Message: msg,
			Details: detail,
		}
	}

	// ErrNetworkError indicates a network-related error
	ErrNetworkError = func(msg string, details ...string) *PluginError {
		var detail string
		if len(details) > 0 {
			detail = details[0]
		}
		return &PluginError{
			Type:    "NetworkError",
			Message: msg,
			Details: detail,
		}
	}

	// ErrUnknown indicates an unknown or unexpected error
	ErrUnknown = func(msg string, details ...string) *PluginError {
		var detail string
		if len(details) > 0 {
			detail = details[0]
		}
		return &PluginError{
			Type:    "Unknown",
			Message: msg,
			Details: detail,
		}
	}
)

// SafeCallCallback wraps a callback execution with error handling
//
// This function provides a safe way to call plugin callbacks with proper
// error handling and logging. It will catch any panics and convert them
// to proper error returns.
//
// Parameters:
//
//	callbackName - Name of the callback for logging
//	callback - The callback function to execute
//
// Returns:
//
//	true if the callback succeeded, false otherwise
//
// Example:
//
//	func (p *MyPlugin) OnEnvironmentCreated(environmentJSON string) bool {
//	    return SafeCallCallback("OnEnvironmentCreated", func() bool {
//	        // Your callback logic here
//	        var env types.Environment
//	        if err := json.Unmarshal([]byte(environmentJSON), &env); err != nil {
//	            return false
//	        }
//
//	        // Process the environment
//	        return p.processEnvironment(&env)
//	    })
//	}
func SafeCallCallback(callbackName string, callback func() bool) bool {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Panic in callback %s: %v\n", callbackName, r)
		}
	}()

	return callback()
}

// LogError logs a plugin error with consistent formatting
func LogError(err *PluginError) {
	fmt.Printf("Plugin Error [%s]: %s\n", err.Type, err.Message)
	if err.Details != "" {
		fmt.Printf("  Details: %s\n", err.Details)
	}
}
