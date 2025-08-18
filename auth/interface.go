package auth

import "github.com/matt953/relm-types-go/types"

// AuthPlugin defines the interface that all authentication plugins must implement
type AuthPlugin interface {
	// CheckUserAccess checks if a user has permission to perform an action on a resource
	CheckUserAccess(userID, resource, action string) (bool, error)

	// CreateUser creates a user in the auth provider
	CreateUser(request types.CreateUserRequest) (*types.CreateUserResult, error)

	// GetUserDetails retrieves user profile/details information
	GetUserDetails(userID string) (*types.UserDetails, error)

	// GetUserDetailsByEmail retrieves user profile/details information by email
	GetUserDetailsByEmail(email string) (*types.UserDetails, error)

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

	// DeleteUser deletes a user from the authentication system
	DeleteUser(userID string) error

	// OAuth Client Management (Required)
	// CreateOAuthClient creates an OAuth client/application
	CreateOAuthClient(request types.CreateOAuthClientRequest) (*types.OAuthClient, error)

	// GetOAuthClient retrieves an OAuth client by client ID
	GetOAuthClient(clientID string) (*types.OAuthClient, error)

	// UpdateOAuthClient updates an OAuth client
	UpdateOAuthClient(clientID string, request types.UpdateOAuthClientRequest) (*types.OAuthClient, error)

	// DeleteOAuthClient deletes an OAuth client
	DeleteOAuthClient(clientID string) error

	// ListOAuthClients lists OAuth clients with pagination
	ListOAuthClients(limit, offset *int) ([]*types.OAuthClient, error)

	// User Authorization Management (Required)
	// ListUserAuthorizedClients lists clients authorized by a user
	ListUserAuthorizedClients(userID string) ([]*types.UserAuthorizedClient, error)

	// RevokeUserClientAuthorization revokes user authorization for a specific client
	RevokeUserClientAuthorization(userID, clientID string) error

	// Cleanup performs any necessary cleanup when the plugin is being unloaded
	// This is optional - plugins can implement this to clean up resources
	Cleanup() error
}

// AuthPluginWithContext extends AuthPlugin with context-aware methods
type AuthPluginWithContext interface {
	AuthPlugin

	// CheckUserAccessWithContext checks user access with additional context
	CheckUserAccessWithContext(userID, resource, action string, context *AuthContext) (bool, error)

	// CreateUserWithContext creates user with additional context
	CreateUserWithContext(request types.CreateUserRequest, context *AuthContext) (*types.CreateUserResult, error)

	// GetUserDetailsByEmailWithContext retrieves user details by email with additional context
	GetUserDetailsByEmailWithContext(email string, context *AuthContext) (*types.UserDetails, error)

	// DeleteUserWithContext deletes a user with additional context
	DeleteUserWithContext(userID string, context *AuthContext) error

	// OAuth Client Management with context
	// CreateOAuthClientWithContext creates OAuth client with additional context
	CreateOAuthClientWithContext(request types.CreateOAuthClientRequest, context *AuthContext) (*types.OAuthClient, error)

	// UpdateOAuthClientWithContext updates OAuth client with additional context
	UpdateOAuthClientWithContext(clientID string, request types.UpdateOAuthClientRequest, context *AuthContext) (*types.OAuthClient, error)

	// DeleteOAuthClientWithContext deletes OAuth client with additional context
	DeleteOAuthClientWithContext(clientID string, context *AuthContext) error

	// RevokeUserClientAuthorizationWithContext revokes user client authorization with additional context
	RevokeUserClientAuthorizationWithContext(userID, clientID string, context *AuthContext) error
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
