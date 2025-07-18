package storage

import "fmt"

// PluginError represents different types of errors that can occur in storage plugins
type PluginError struct {
	Type    ErrorType
	Message string
}

// ErrorType represents the category of error
type ErrorType int

const (
	InvalidInputError ErrorType = iota
	StorageErrorType
	NetworkErrorType  
	ConfigurationErrorType
	UnknownErrorType
)

func (e *PluginError) Error() string {
	switch e.Type {
	case InvalidInputError:
		return fmt.Sprintf("Invalid input: %s", e.Message)
	case StorageErrorType:
		return fmt.Sprintf("Storage error: %s", e.Message)
	case NetworkErrorType:
		return fmt.Sprintf("Network error: %s", e.Message)
	case ConfigurationErrorType:
		return fmt.Sprintf("Configuration error: %s", e.Message)
	case UnknownErrorType:
		return fmt.Sprintf("Unknown error: %s", e.Message)
	default:
		return fmt.Sprintf("Unknown error: %s", e.Message)
	}
}

// NewInvalidInputError creates a new invalid input error
func NewInvalidInputError(message string) *PluginError {
	return &PluginError{
		Type:    InvalidInputError,
		Message: message,
	}
}

// NewStorageError creates a new storage error
func NewStorageError(message string) *PluginError {
	return &PluginError{
		Type:    StorageErrorType,
		Message: message,
	}
}

// NewNetworkError creates a new network error
func NewNetworkError(message string) *PluginError {
	return &PluginError{
		Type:    NetworkErrorType,
		Message: message,
	}
}

// NewConfigurationError creates a new configuration error
func NewConfigurationError(message string) *PluginError {
	return &PluginError{
		Type:    ConfigurationErrorType,
		Message: message,
	}
}

// NewUnknownError creates a new unknown error
func NewUnknownError(message string) *PluginError {
	return &PluginError{
		Type:    UnknownErrorType,
		Message: message,
	}
}