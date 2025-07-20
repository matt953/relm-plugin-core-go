package config

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/matt953/relm-types-go/types"
)

// Re-export types for backward compatibility
type PluginConfigItem = types.PluginConfigItem
type PluginDefinition = types.PluginDefinition
type RealmConfig = types.RealmConfig

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

// GetConfigValue gets a configuration value from the global config
func GetConfigValue(key string) (string, bool) {
	configMutex.RLock()
	defer configMutex.RUnlock()
	
	// Check plugin_config section
	if pluginConfig, ok := globalConfig["plugin_config"].(map[string]interface{}); ok {
		if value, exists := pluginConfig[key]; exists {
			if strVal, ok := value.(string); ok {
				return strVal, true
			}
		}
	}
	
	// Check server_config section
	if serverConfig, ok := globalConfig["server_config"].(map[string]interface{}); ok {
		if value, exists := serverConfig[key]; exists {
			if strVal, ok := value.(string); ok {
				return strVal, true
			}
		}
	}
	
	return "", false
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

// GetServerConfigValue gets a configuration value specifically from the server config section
func GetServerConfigValue(key string) (string, bool) {
	configMutex.RLock()
	defer configMutex.RUnlock()
	
	if serverConfig, ok := globalConfig["server_config"].(map[string]interface{}); ok {
		if value, exists := serverConfig[key]; exists {
			if strVal, ok := value.(string); ok {
				return strVal, true
			}
		}
	}
	
	return "", false
}

// LoadConfigAndSetEnvVars works with the global config that was passed via initialize_with_config
// For backward compatibility, this function now returns success if global config is available
func LoadConfigAndSetEnvVars(pluginType string) error {
	configMutex.RLock()
	defer configMutex.RUnlock()
	
	if globalConfig == nil {
		// Fallback to old method if global config is not set
		return types.LoadConfigAndSetEnvVars(pluginType)
	}
	
	// Global config is available, no need to load from files
	return nil
}

// MustLoadConfigAndSetEnvVars is like LoadConfigAndSetEnvVars but panics on error
func MustLoadConfigAndSetEnvVars(pluginType string) {
	if err := LoadConfigAndSetEnvVars(pluginType); err != nil {
		panic(err)
	}
}

// LoadGeneralPluginConfig works with the global config for backward compatibility
func LoadGeneralPluginConfig(pluginName string) error {
	configMutex.RLock()
	defer configMutex.RUnlock()
	
	if globalConfig == nil {
		// Fallback to old method if global config is not set
		return types.LoadGeneralPluginConfig(pluginName)
	}
	
	// Global config is available, no need to load from files
	return nil
}

// LoadAllGeneralPluginConfigs works with the global config for backward compatibility
func LoadAllGeneralPluginConfigs() error {
	configMutex.RLock()
	defer configMutex.RUnlock()
	
	if globalConfig == nil {
		// Fallback to old method if global config is not set
		return types.LoadAllGeneralPluginConfigs()
	}
	
	// Global config is available, no need to load from files
	return nil
}

// InMemoryConfigProvider provides access to in-memory configuration
type InMemoryConfigProvider struct {
	config *types.InMemoryConfig
	useGlobalConfig bool
}

// NewInMemoryConfigProvider creates a new in-memory configuration provider
func NewInMemoryConfigProvider(config *types.InMemoryConfig) *InMemoryConfigProvider {
	return &InMemoryConfigProvider{config: config, useGlobalConfig: false}
}

// NewInMemoryConfigProviderFromGlobal creates a new provider from the global config
func NewInMemoryConfigProviderFromGlobal() *InMemoryConfigProvider {
	return &InMemoryConfigProvider{config: nil, useGlobalConfig: true}
}

// GetGlobalConfigProvider returns a config provider that uses the global config
// This is the recommended way for plugins to access configuration
func GetGlobalConfigProvider() *InMemoryConfigProvider {
	return NewInMemoryConfigProviderFromGlobal()
}

// GetPluginConfig gets a configuration value for a specific plugin type and key
func (p *InMemoryConfigProvider) GetPluginConfig(pluginType, key string) (string, bool) {
	if p.useGlobalConfig {
		envVar := fmt.Sprintf("%s_%s", strings.ToUpper(pluginType), strings.ToUpper(key))
		return GetPluginConfigValue(envVar)
	}
	
	envVar := fmt.Sprintf("%s_%s", strings.ToUpper(pluginType), strings.ToUpper(key))
	if pluginConfig := p.config.PluginConfigs[strings.ToUpper(pluginType)]; pluginConfig != nil {
		value, exists := pluginConfig[envVar]
		return value, exists
	}
	return "", false
}

// GetConfigValue gets a configuration value from a specific section
func (p *InMemoryConfigProvider) GetConfigValue(section, key string) (string, bool) {
	if p.useGlobalConfig {
		return GetConfigValue(key)
	}
	return p.config.GetConfigValue(section, key)
}

// GetPluginConfigMap gets all configuration for a specific plugin
func (p *InMemoryConfigProvider) GetPluginConfigMap(pluginKey string) map[string]string {
	if p.useGlobalConfig {
		result := make(map[string]string)
		configMutex.RLock()
		defer configMutex.RUnlock()
		
		if pluginConfig, ok := globalConfig["plugin_config"].(map[string]interface{}); ok {
			for k, v := range pluginConfig {
				if strVal, ok := v.(string); ok {
					result[k] = strVal
				}
			}
		}
		return result
	}
	return p.config.GetPluginConfig(pluginKey)
}

// GetBool gets a boolean configuration value
func (p *InMemoryConfigProvider) GetBool(key string) bool {
	if value, exists := p.GetConfigByKey(key); exists {
		v := strings.ToLower(value)
		return v == "true" || v == "1" || v == "yes" || v == "on"
	}
	return false
}

// GetInt gets an integer configuration value
func (p *InMemoryConfigProvider) GetInt(key string) (int64, error) {
	if value, exists := p.GetConfigByKey(key); exists {
		return strconv.ParseInt(value, 10, 64)
	}
	return 0, fmt.Errorf("configuration '%s' not found", key)
}

// GetOrDefault gets a configuration value with a default
func (p *InMemoryConfigProvider) GetOrDefault(key, defaultValue string) string {
	if value, exists := p.GetConfigByKey(key); exists {
		return value
	}
	return defaultValue
}

// GetConfigByKey searches all sections for a configuration key
func (p *InMemoryConfigProvider) GetConfigByKey(key string) (string, bool) {
	if p.useGlobalConfig {
		return GetConfigValue(key)
	}
	
	// Check server config
	if value, exists := p.config.ServerConfig[key]; exists {
		return value, true
	}

	// Check database config
	if value, exists := p.config.DatabaseConfig[key]; exists {
		return value, true
	}

	// Check storage config
	if value, exists := p.config.StorageConfig[key]; exists {
		return value, true
	}

	// Check JWT config
	if value, exists := p.config.JWTConfig[key]; exists {
		return value, true
	}

	// Check plugin configs
	for _, pluginConfig := range p.config.PluginConfigs {
		if value, exists := pluginConfig[key]; exists {
			return value, true
		}
	}

	return "", false
}

// LoadConfigAndBuildInMemory loads the JSON config and returns in-memory config provider
// If global config is available, it returns a provider that uses the global config
func LoadConfigAndBuildInMemory() (*InMemoryConfigProvider, error) {
	configMutex.RLock()
	hasGlobalConfig := globalConfig != nil
	configMutex.RUnlock()
	
	if hasGlobalConfig {
		return GetGlobalConfigProvider(), nil
	}
	
	config, err := types.LoadConfigAndBuildInMemory()
	if err != nil {
		return nil, err
	}
	return NewInMemoryConfigProvider(config), nil
}

// PluginConfigInterface defines the interface for plugin configuration
type PluginConfigInterface interface {
	GetConfig(provider *InMemoryConfigProvider, key string) (string, error)
	GetConfigOptional(provider *InMemoryConfigProvider, key string) (string, bool)
	GetConfigBool(provider *InMemoryConfigProvider, key string) bool
	GetConfigInt(provider *InMemoryConfigProvider, key string) (int64, error)
	GetConfigOrDefault(provider *InMemoryConfigProvider, key, defaultValue string) string
	PluginType() string
}

// BasePluginConfig provides base implementation for plugin configuration
type BasePluginConfig struct {
	pluginType string
}

// NewBasePluginConfig creates a new base plugin configuration
func NewBasePluginConfig(pluginType string) *BasePluginConfig {
	return &BasePluginConfig{pluginType: strings.ToUpper(pluginType)}
}

// PluginType returns the plugin type
func (c *BasePluginConfig) PluginType() string {
	return c.pluginType
}

// GetConfig gets a required configuration value
func (c *BasePluginConfig) GetConfig(provider *InMemoryConfigProvider, key string) (string, error) {
	if value, exists := provider.GetPluginConfig(c.pluginType, key); exists {
		return value, nil
	}
	return "", fmt.Errorf("configuration '%s' not found for plugin type '%s'", key, c.pluginType)
}

// GetConfigOptional gets an optional configuration value
func (c *BasePluginConfig) GetConfigOptional(provider *InMemoryConfigProvider, key string) (string, bool) {
	return provider.GetPluginConfig(c.pluginType, key)
}

// GetConfigBool gets a boolean configuration value
func (c *BasePluginConfig) GetConfigBool(provider *InMemoryConfigProvider, key string) bool {
	envVar := fmt.Sprintf("%s_%s", c.pluginType, strings.ToUpper(key))
	return provider.GetBool(envVar)
}

// GetConfigInt gets an integer configuration value
func (c *BasePluginConfig) GetConfigInt(provider *InMemoryConfigProvider, key string) (int64, error) {
	envVar := fmt.Sprintf("%s_%s", c.pluginType, strings.ToUpper(key))
	return provider.GetInt(envVar)
}

// GetConfigOrDefault gets a configuration value with a default
func (c *BasePluginConfig) GetConfigOrDefault(provider *InMemoryConfigProvider, key, defaultValue string) string {
	envVar := fmt.Sprintf("%s_%s", c.pluginType, strings.ToUpper(key))
	return provider.GetOrDefault(envVar, defaultValue)
}

// Predefined plugin configurations
var (
	StorageConfig        = NewBasePluginConfig("STORAGE")
	AuthenticationConfig = NewBasePluginConfig("AUTHENTICATION")
	ProcessingConfig     = NewBasePluginConfig("PROCESSING")
	AnalyticsConfig      = NewBasePluginConfig("ANALYTICS")
)

// NewGeneralConfig creates a configuration for general plugins
func NewGeneralConfig(pluginName string) *BasePluginConfig {
	return NewBasePluginConfig(strings.ToUpper(strings.ReplaceAll(pluginName, "-", "_")))
}

// Convenience functions for plugins to easily access configuration

// GetBool gets a boolean configuration value from the global config
func GetBool(key string) bool {
	if value, exists := GetConfigValue(key); exists {
		v := strings.ToLower(value)
		return v == "true" || v == "1" || v == "yes" || v == "on"
	}
	return false
}

// GetInt gets an integer configuration value from the global config
func GetInt(key string) (int64, error) {
	if value, exists := GetConfigValue(key); exists {
		return strconv.ParseInt(value, 10, 64)
	}
	return 0, fmt.Errorf("configuration '%s' not found", key)
}

// GetOrDefault gets a configuration value with a default from the global config
func GetOrDefault(key, defaultValue string) string {
	if value, exists := GetConfigValue(key); exists {
		return value
	}
	return defaultValue
}

// GetPluginBool gets a boolean configuration value from the plugin config section
func GetPluginBool(key string) bool {
	if value, exists := GetPluginConfigValue(key); exists {
		v := strings.ToLower(value)
		return v == "true" || v == "1" || v == "yes" || v == "on"
	}
	return false
}

// GetPluginInt gets an integer configuration value from the plugin config section
func GetPluginInt(key string) (int64, error) {
	if value, exists := GetPluginConfigValue(key); exists {
		return strconv.ParseInt(value, 10, 64)
	}
	return 0, fmt.Errorf("plugin configuration '%s' not found", key)
}

// GetPluginOrDefault gets a plugin configuration value with a default
func GetPluginOrDefault(key, defaultValue string) string {
	if value, exists := GetPluginConfigValue(key); exists {
		return value
	}
	return defaultValue
}

// GetServerBool gets a boolean configuration value from the server config section
func GetServerBool(key string) bool {
	if value, exists := GetServerConfigValue(key); exists {
		v := strings.ToLower(value)
		return v == "true" || v == "1" || v == "yes" || v == "on"
	}
	return false
}

// GetServerInt gets an integer configuration value from the server config section
func GetServerInt(key string) (int64, error) {
	if value, exists := GetServerConfigValue(key); exists {
		return strconv.ParseInt(value, 10, 64)
	}
	return 0, fmt.Errorf("server configuration '%s' not found", key)
}

// GetServerOrDefault gets a server configuration value with a default
func GetServerOrDefault(key, defaultValue string) string {
	if value, exists := GetServerConfigValue(key); exists {
		return value
	}
	return defaultValue
}
