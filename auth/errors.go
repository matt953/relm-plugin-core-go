package auth

import "fmt"

// PluginError represents different types of errors that can occur in auth plugins
type PluginError struct {
	Type    ErrorType
	Message string
}

// ErrorType represents the category of error
type ErrorType int

const (
	InvalidInputError ErrorType = iota
	AuthenticationErrorType
	AuthorizationErrorType
	UserNotFoundErrorType
	PermissionDeniedErrorType
	NetworkErrorType
	ConfigurationErrorType
	InitializationErrorType
	SerializationErrorType
	OperationFailedErrorType
	UnknownErrorType
)

func (e *PluginError) Error() string {
	switch e.Type {
	case InvalidInputError:
		return fmt.Sprintf("Invalid input: %s", e.Message)
	case AuthenticationErrorType:
		return fmt.Sprintf("Authentication error: %s", e.Message)
	case AuthorizationErrorType:
		return fmt.Sprintf("Authorization error: %s", e.Message)
	case UserNotFoundErrorType:
		return fmt.Sprintf("User not found: %s", e.Message)
	case PermissionDeniedErrorType:
		return fmt.Sprintf("Permission denied: %s", e.Message)
	case NetworkErrorType:
		return fmt.Sprintf("Network error: %s", e.Message)
	case ConfigurationErrorType:
		return fmt.Sprintf("Configuration error: %s", e.Message)
	case InitializationErrorType:
		return fmt.Sprintf("Initialization error: %s", e.Message)
	case SerializationErrorType:
		return fmt.Sprintf("Serialization error: %s", e.Message)
	case OperationFailedErrorType:
		return fmt.Sprintf("Operation failed: %s", e.Message)
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

// NewAuthenticationError creates a new authentication error
func NewAuthenticationError(message string) *PluginError {
	return &PluginError{
		Type:    AuthenticationErrorType,
		Message: message,
	}
}

// NewAuthorizationError creates a new authorization error
func NewAuthorizationError(message string) *PluginError {
	return &PluginError{
		Type:    AuthorizationErrorType,
		Message: message,
	}
}

// NewUserNotFoundError creates a new user not found error
func NewUserNotFoundError(message string) *PluginError {
	return &PluginError{
		Type:    UserNotFoundErrorType,
		Message: message,
	}
}

// NewPermissionDeniedError creates a new permission denied error
func NewPermissionDeniedError(message string) *PluginError {
	return &PluginError{
		Type:    PermissionDeniedErrorType,
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

// NewInitializationError creates a new initialization error
func NewInitializationError(message string) *PluginError {
	return &PluginError{
		Type:    InitializationErrorType,
		Message: message,
	}
}

// NewSerializationError creates a new serialization error
func NewSerializationError(message string) *PluginError {
	return &PluginError{
		Type:    SerializationErrorType,
		Message: message,
	}
}

// NewOperationFailedError creates a new operation failed error
func NewOperationFailedError(message string) *PluginError {
	return &PluginError{
		Type:    OperationFailedErrorType,
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