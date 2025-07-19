package storage

import (
	"github.com/matt953/relm-plugin-core-go/config"
)

// StoragePlugin defines the interface that all storage plugins must implement
type StoragePlugin interface {
	// StoreFile stores data at the specified path with optional content type
	StoreFile(path string, data []byte, contentType *string) error

	// RetrieveFile retrieves file data from the specified path
	RetrieveFile(path string) ([]byte, error)

	// DeleteFile deletes the file at the specified path
	DeleteFile(path string) error

	// FileExists checks if a file exists at the specified path
	FileExists(path string) bool

	// GenerateURL generates a public URL for accessing the stored file
	// Returns nil if URL generation is not supported
	GenerateURL(path string, baseURL string) *string

	// ProviderName returns a human-readable name for this storage provider
	ProviderName() string

	// Cleanup performs any necessary cleanup when the plugin is being unloaded
	// This is optional - plugins can implement this to clean up resources
	Cleanup() error
}

// Global variable to hold the registered plugin instance
var registeredPlugin StoragePlugin

// Plugin initializer callback function type
type PluginInitializer func() (StoragePlugin, error)

// Global plugin initializer callback
var pluginInitializer PluginInitializer

// RegisterPlugin registers a storage plugin for FFI export
// This must be called from the plugin's main function
func RegisterPlugin(plugin StoragePlugin) {
	registeredPlugin = plugin
}

// SetPluginInitializer sets the callback function for lazy plugin initialization
// This should be called from the plugin's main package
func SetPluginInitializer(initializer PluginInitializer) {
	pluginInitializer = initializer
}

// GetRegisteredPlugin returns the currently registered plugin
// Used internally by FFI functions
func GetRegisteredPlugin() StoragePlugin {
	// If no plugin is registered but we have an initializer, try to initialize
	if registeredPlugin == nil && pluginInitializer != nil {
		// Load configuration and set environment variables for storage plugins
		// This happens automatically before plugin initialization
		config.LoadConfigAndSetEnvVars("storage")
		
		plugin, err := pluginInitializer()
		if err == nil && plugin != nil {
			RegisterPlugin(plugin)
		}
	}
	return registeredPlugin
}
