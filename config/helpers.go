package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// GetConfig gets a plugin configuration value from environment variables
//
// This function automatically prepends the plugin type prefix to the variable name.
// For example, if you're in a storage plugin and call GetConfig("STORAGE", "bucket"),
// it will look for the environment variable STORAGE_BUCKET.
//
// Parameters:
//   - pluginType: The type of plugin (e.g., "STORAGE", "AUTHENTICATION", "PROCESSING", "ANALYTICS", or a custom name for general plugins)
//   - key: The configuration key (will be converted to uppercase)
//
// Returns:
//   - The configuration value if found
//   - An error if the configuration is not found
//
// Example:
//
//	// In a storage plugin
//	bucket, err := config.GetConfig("STORAGE", "bucket")
//	region, err := config.GetConfig("STORAGE", "region")
//
//	// In a general plugin (using the plugin name)
//	webhookURL, err := config.GetConfig("NOTIFICATION_PLUGIN", "webhook_url")
func GetConfig(pluginType, key string) (string, error) {
	envVar := fmt.Sprintf("%s_%s", strings.ToUpper(pluginType), strings.ToUpper(key))
	value := os.Getenv(envVar)
	if value == "" {
		return "", fmt.Errorf("configuration '%s' not found (looking for environment variable: %s)", key, envVar)
	}
	return value, nil
}

// GetConfigOptional gets an optional plugin configuration value from environment variables
//
// Similar to GetConfig but returns an empty string if the configuration is not found
// instead of an error.
//
// Parameters:
//   - pluginType: The type of plugin
//   - key: The configuration key (will be converted to uppercase)
//
// Returns:
//   - The configuration value if found, empty string otherwise
//
// Example:
//
//	// In a storage plugin
//	customEndpoint := config.GetConfigOptional("STORAGE", "custom_endpoint")
//	if customEndpoint != "" {
//	    // Use custom endpoint
//	}
func GetConfigOptional(pluginType, key string) string {
	envVar := fmt.Sprintf("%s_%s", strings.ToUpper(pluginType), strings.ToUpper(key))
	return os.Getenv(envVar)
}

// GetConfigBool gets a boolean configuration value
//
// Recognizes "true", "1", "yes", "on" as true (case-insensitive).
// All other values (including missing) are considered false.
//
// Parameters:
//   - pluginType: The type of plugin
//   - key: The configuration key
//
// Returns:
//   - true if the value is a recognized truthy value
//   - false otherwise
//
// Example:
//
//	// In a storage plugin
//	forcePathStyle := config.GetConfigBool("STORAGE", "force_path_style")
func GetConfigBool(pluginType, key string) bool {
	value := strings.ToLower(GetConfigOptional(pluginType, key))
	return value == "true" || value == "1" || value == "yes" || value == "on"
}

// GetConfigInt gets an integer configuration value
//
// Parameters:
//   - pluginType: The type of plugin
//   - key: The configuration key
//
// Returns:
//   - The parsed integer value
//   - An error if the configuration is not found or cannot be parsed
//
// Example:
//
//	// In a processing plugin
//	maxWorkers, err := config.GetConfigInt("PROCESSING", "max_workers")
func GetConfigInt(pluginType, key string) (int, error) {
	value, err := GetConfig(pluginType, key)
	if err != nil {
		return 0, err
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("configuration '%s' is not a valid integer: %v", key, err)
	}
	
	return intValue, nil
}

// GetConfigIntOptional gets an optional integer configuration value
//
// Returns the default value if the configuration is not found or cannot be parsed.
//
// Parameters:
//   - pluginType: The type of plugin
//   - key: The configuration key
//   - defaultValue: The default value to use if not found or invalid
//
// Returns:
//   - The configuration value if found and valid, otherwise the default value
func GetConfigIntOptional(pluginType, key string, defaultValue int) int {
	value := GetConfigOptional(pluginType, key)
	if value == "" {
		return defaultValue
	}
	
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	
	return intValue
}

// GetConfigOrDefault gets a configuration value with a default
//
// Parameters:
//   - pluginType: The type of plugin
//   - key: The configuration key
//   - defaultValue: The default value to use if not found
//
// Returns:
//   - The configuration value if found, otherwise the default value
//
// Example:
//
//	// In an analytics plugin
//	batchSize := config.GetConfigOrDefault("ANALYTICS", "batch_size", "100")
func GetConfigOrDefault(pluginType, key, defaultValue string) string {
	value := GetConfigOptional(pluginType, key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Config provides a convenient interface for accessing plugin configuration
type Config struct {
	pluginType string
}

// NewConfig creates a new Config instance for the given plugin type
//
// Example:
//
//	// For storage plugins
//	cfg := config.NewConfig("STORAGE")
//	bucket, err := cfg.Get("bucket")
//
//	// For general plugins
//	cfg := config.NewConfig("NOTIFICATION_PLUGIN")
//	webhookURL, err := cfg.Get("webhook_url")
func NewConfig(pluginType string) *Config {
	return &Config{
		pluginType: strings.ToUpper(strings.ReplaceAll(pluginType, "-", "_")),
	}
}

// Get gets a required configuration value
func (c *Config) Get(key string) (string, error) {
	return GetConfig(c.pluginType, key)
}

// GetOptional gets an optional configuration value
func (c *Config) GetOptional(key string) string {
	return GetConfigOptional(c.pluginType, key)
}

// GetBool gets a boolean configuration value
func (c *Config) GetBool(key string) bool {
	return GetConfigBool(c.pluginType, key)
}

// GetInt gets an integer configuration value
func (c *Config) GetInt(key string) (int, error) {
	return GetConfigInt(c.pluginType, key)
}

// GetIntOptional gets an optional integer configuration value with a default
func (c *Config) GetIntOptional(key string, defaultValue int) int {
	return GetConfigIntOptional(c.pluginType, key, defaultValue)
}

// GetOrDefault gets a configuration value with a default
func (c *Config) GetOrDefault(key, defaultValue string) string {
	return GetConfigOrDefault(c.pluginType, key, defaultValue)
}

// Predefined config instances for common plugin types
var (
	// Storage provides configuration access for storage plugins
	Storage = NewConfig("STORAGE")
	
	// Authentication provides configuration access for authentication plugins
	Authentication = NewConfig("AUTHENTICATION")
	
	// Processing provides configuration access for processing plugins
	Processing = NewConfig("PROCESSING")
	
	// Analytics provides configuration access for analytics plugins
	Analytics = NewConfig("ANALYTICS")
)