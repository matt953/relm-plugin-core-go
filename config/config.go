package config

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Global configuration management
var (
	globalConfig map[string]interface{}
	configMutex  sync.RWMutex
)

// SetConfigFromJSON sets the global config from JSON string
// This is called by the FFI initialize_with_config function
func SetConfigFromJSON(configJSON string) error {
	configMutex.Lock()
	defer configMutex.Unlock()

	if configJSON == "" {
		return fmt.Errorf("empty config JSON")
	}

	if err := json.Unmarshal([]byte(configJSON), &globalConfig); err != nil {
		return fmt.Errorf("failed to parse config JSON: %v", err)
	}

	return nil
}

// GetGlobalConfig returns a copy of the global config
func GetGlobalConfig() map[string]interface{} {
	configMutex.RLock()
	defer configMutex.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]interface{})
	for k, v := range globalConfig {
		result[k] = v
	}
	return result
}

// GetPluginConfigValue gets a configuration value specifically from the plugin config section
func GetPluginConfigValue(key string) (string, bool) {
	configMutex.RLock()
	defer configMutex.RUnlock()

	if pluginConfig, ok := globalConfig["plugin_config"].(map[string]interface{}); ok {
		if value, exists := pluginConfig[key]; exists {
			if strVal, ok := value.(string); ok {
				return strVal, true
			}
		}
	}

	return "", false
}

// GetPluginBool gets a boolean configuration value from the plugin config section
func GetPluginBool(key string) bool {
	if value, exists := GetPluginConfigValue(key); exists {
		v := strings.ToLower(value)
		return v == "true" || v == "1" || v == "yes" || v == "on"
	}
	return false
}

// GetPluginOrDefault gets a plugin-specific configuration value with a default
func GetPluginOrDefault(key, defaultValue string) string {
	if value, exists := GetPluginConfigValue(key); exists {
		return value
	}
	return defaultValue
}