package config

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/matt953/relm-types-go/types"
)

// Re-export types for backward compatibility
type PluginConfigItem = types.PluginConfigItem
type PluginDefinition = types.PluginDefinition
type RealmConfig = types.RealmConfig

// LoadConfigAndSetEnvVars loads the JSON config and sets environment variables for the specified plugin type
func LoadConfigAndSetEnvVars(pluginType string) error {
	return types.LoadConfigAndSetEnvVars(pluginType)
}

// MustLoadConfigAndSetEnvVars is like LoadConfigAndSetEnvVars but panics on error
func MustLoadConfigAndSetEnvVars(pluginType string) {
	types.MustLoadConfigAndSetEnvVars(pluginType)
}

// LoadGeneralPluginConfig loads configuration for a specific general plugin by finding it in the JSON config
// and setting environment variables with the extracted plugin name
func LoadGeneralPluginConfig(pluginName string) error {
	return types.LoadGeneralPluginConfig(pluginName)
}

// LoadAllGeneralPluginConfigs loads configuration for all general plugins in the JSON config
// This is useful when a Go general plugin doesn't know its specific name but wants to load
// environment variables for all general plugins
func LoadAllGeneralPluginConfigs() error {
	return types.LoadAllGeneralPluginConfigs()
}

// InMemoryConfigProvider provides access to in-memory configuration
type InMemoryConfigProvider struct {
	config *types.InMemoryConfig
}

// NewInMemoryConfigProvider creates a new in-memory configuration provider
func NewInMemoryConfigProvider(config *types.InMemoryConfig) *InMemoryConfigProvider {
	return &InMemoryConfigProvider{config: config}
}

// GetPluginConfig gets a configuration value for a specific plugin type and key
func (p *InMemoryConfigProvider) GetPluginConfig(pluginType, key string) (string, bool) {
	envVar := fmt.Sprintf("%s_%s", strings.ToUpper(pluginType), strings.ToUpper(key))
	if pluginConfig := p.config.PluginConfigs[strings.ToUpper(pluginType)]; pluginConfig != nil {
		value, exists := pluginConfig[envVar]
		return value, exists
	}
	return "", false
}

// GetConfigValue gets a configuration value from a specific section
func (p *InMemoryConfigProvider) GetConfigValue(section, key string) (string, bool) {
	return p.config.GetConfigValue(section, key)
}

// GetPluginConfigMap gets all configuration for a specific plugin
func (p *InMemoryConfigProvider) GetPluginConfigMap(pluginKey string) map[string]string {
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
func LoadConfigAndBuildInMemory() (*InMemoryConfigProvider, error) {
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
