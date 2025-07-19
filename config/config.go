package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// PluginConfigItem represents a single configuration item
type PluginConfigItem struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// PluginDefinition represents a single plugin definition
type PluginDefinition struct {
	Type    string             `yaml:"type"`
	Enabled bool               `yaml:"enabled"`
	Path    string             `yaml:"path"`
	Config  []PluginConfigItem `yaml:"config"`
}

// RealmConfig represents the full realm configuration
type RealmConfig struct {
	Plugins []PluginDefinition `yaml:"plugins"` // Array-based format (all plugins are FFI)
}

// LoadConfigAndSetEnvVars loads the YAML config and sets environment variables for the specified plugin type
func LoadConfigAndSetEnvVars(pluginType string) error {
	// Try to get config file path from environment variable first
	configPath := os.Getenv("REALM_CONFIG_FILE")
	if configPath == "" {
		// Fall back to common locations
		possiblePaths := []string{
			"realm.yaml",
			"config.yaml",
			"../realm.yaml",
			"../../realm.yaml",
			"../../../realm.yaml",
		}
		
		for _, path := range possiblePaths {
			if _, err := os.Stat(path); err == nil {
				configPath = path
				break
			}
		}
	}

	if configPath == "" {
		// No config file found, which is fine - plugins can still work with direct env vars
		return nil
	}

	// Make path absolute
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for config file %s: %v", configPath, err)
	}

	// Read the config file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", absPath, err)
	}

	// Parse YAML
	var config RealmConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse YAML config: %v", err)
	}

	// Set environment variables for the specific plugin type from array format
	if config.Plugins != nil {
		for _, plugin := range config.Plugins {
			if strings.EqualFold(plugin.Type, pluginType) {
				for _, configItem := range plugin.Config {
					envVarName := fmt.Sprintf("%s_%s", strings.ToUpper(pluginType), strings.ToUpper(configItem.Name))
					// Only set if not already set (allow env vars to override config file)
					if os.Getenv(envVarName) == "" {
						os.Setenv(envVarName, configItem.Value)
					}
				}
				break // Found the plugin, no need to continue
			}
		}
	}

	return nil
}

// MustLoadConfigAndSetEnvVars is like LoadConfigAndSetEnvVars but panics on error
func MustLoadConfigAndSetEnvVars(pluginType string) {
	if err := LoadConfigAndSetEnvVars(pluginType); err != nil {
		panic(fmt.Sprintf("Failed to load config for plugin type %s: %v", pluginType, err))
	}
}