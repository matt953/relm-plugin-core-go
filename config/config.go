package config

import (
	"github.com/matt953/relm-types-go/types"
)

// Re-export types for backward compatibility
type PluginConfigItem = types.PluginConfigItem
type PluginDefinition = types.PluginDefinition
type RealmConfig = types.RealmConfig

// LoadConfigAndSetEnvVars loads the YAML config and sets environment variables for the specified plugin type
func LoadConfigAndSetEnvVars(pluginType string) error {
	return types.LoadConfigAndSetEnvVars(pluginType)
}

// MustLoadConfigAndSetEnvVars is like LoadConfigAndSetEnvVars but panics on error
func MustLoadConfigAndSetEnvVars(pluginType string) {
	types.MustLoadConfigAndSetEnvVars(pluginType)
}

// LoadGeneralPluginConfig loads configuration for a specific general plugin by finding it in the YAML config
// and setting environment variables with the extracted plugin name
func LoadGeneralPluginConfig(pluginName string) error {
	return types.LoadGeneralPluginConfig(pluginName)
}

// LoadAllGeneralPluginConfigs loads configuration for all general plugins in the YAML config
// This is useful when a Go general plugin doesn't know its specific name but wants to load
// environment variables for all general plugins
func LoadAllGeneralPluginConfigs() error {
	return types.LoadAllGeneralPluginConfigs()
}
