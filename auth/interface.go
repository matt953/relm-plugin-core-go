package auth

import "github.com/matt953/relm-types-go/types"

// AuthPlugin defines the interface that all authentication plugins must implement
type AuthPlugin interface {
	// CheckUserAccess checks if a user has permission to perform an action on a resource
	CheckUserAccess(userID, resource, action string) (bool, error)

	// GetUserPermissions retrieves detailed permissions for a user
	GetUserPermissions(userID string) (*types.UserPermissions, error)

	// GetUserDetails retrieves user profile/details information
	GetUserDetails(userID string) (*types.UserDetails, error)

	// GetPluginInfo returns plugin information and capabilities
	GetPluginInfo() (*types.AuthPluginInfo, error)

	// ProviderName returns a human-readable name for this auth provider
	ProviderName() string

	// Initialize initializes the plugin with configuration (JSON string)
	Initialize(configJSON *string) error

	// HealthCheck checks if the plugin is healthy and can respond to requests
	HealthCheck() bool

	// ValidateUser checks if a user exists in the authentication system
	ValidateUser(userID string) bool

	// GetUserGroups returns the groups/roles for a user (optional, can return empty slice)
	GetUserGroups(userID string) ([]string, error)

	// SearchUsers searches for users matching a query (optional, can return empty slice)
	SearchUsers(query string, limit int) ([]*types.UserDetails, error)

	// Cleanup performs any necessary cleanup when the plugin is being unloaded
	// This is optional - plugins can implement this to clean up resources
	Cleanup() error
}

// AuthPluginWithContext extends AuthPlugin with context-aware methods
type AuthPluginWithContext interface {
	AuthPlugin

	// CheckUserAccessWithContext checks user access with additional context
	CheckUserAccessWithContext(userID, resource, action string, context *AuthContext) (bool, error)

	// GetUserPermissionsWithContext gets user permissions with additional context
	GetUserPermissionsWithContext(userID string, context *AuthContext) (*types.UserPermissions, error)
}

// AuthContext provides additional context for authentication operations
type AuthContext struct {
	RequestID      *string                `json:"request_id,omitempty"`
	ClientIP       *string                `json:"client_ip,omitempty"`
	UserAgent      *string                `json:"user_agent,omitempty"`
	AdditionalData map[string]interface{} `json:"additional_data,omitempty"`
}

// NewAuthContext creates a new AuthContext
func NewAuthContext() *AuthContext {
	return &AuthContext{
		AdditionalData: make(map[string]interface{}),
	}
}

// WithRequestID sets the request ID
func (ac *AuthContext) WithRequestID(requestID string) *AuthContext {
	ac.RequestID = &requestID
	return ac
}

// WithClientIP sets the client IP
func (ac *AuthContext) WithClientIP(clientIP string) *AuthContext {
	ac.ClientIP = &clientIP
	return ac
}

// WithUserAgent sets the user agent
func (ac *AuthContext) WithUserAgent(userAgent string) *AuthContext {
	ac.UserAgent = &userAgent
	return ac
}

// WithData adds additional data
func (ac *AuthContext) WithData(key string, value interface{}) *AuthContext {
	ac.AdditionalData[key] = value
	return ac
}

// Global variable to hold the registered plugin instance
var registeredPlugin AuthPlugin

// Plugin initializer callback function type
type PluginInitializer func() (AuthPlugin, error)

// Global plugin initializer callback
var pluginInitializer PluginInitializer

// RegisterPlugin registers a storage plugin for FFI export
// This must be called from the plugin's main function
func RegisterPlugin(plugin AuthPlugin) {
	registeredPlugin = plugin
}

// SetPluginInitializer sets the callback function for lazy plugin initialization
// This should be called from the plugin's main package
func SetPluginInitializer(initializer PluginInitializer) {
	pluginInitializer = initializer
}

// GetRegisteredPlugin returns the currently registered plugin
// Used internally by FFI functions
func GetRegisteredPlugin() AuthPlugin {
	// If no plugin is registered but we have an initializer, try to initialize
	if registeredPlugin == nil && pluginInitializer != nil {
		plugin, err := pluginInitializer()
		if err == nil && plugin != nil {
			RegisterPlugin(plugin)
		}
	}
	return registeredPlugin
}
