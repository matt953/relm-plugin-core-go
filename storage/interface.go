package storage

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
}

// Global variable to hold the registered plugin instance
var registeredPlugin StoragePlugin

// RegisterPlugin registers a storage plugin for FFI export
// This must be called from the plugin's main function
func RegisterPlugin(plugin StoragePlugin) {
	registeredPlugin = plugin
}

// GetRegisteredPlugin returns the currently registered plugin
// Used internally by FFI functions
func GetRegisteredPlugin() StoragePlugin {
	return registeredPlugin
}