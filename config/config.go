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
					var envVarName string
					if strings.ToLower(pluginType) == "general" {
						// For general plugins, extract plugin name from path like Rust does
						pluginName := "GENERAL_UNKNOWN"
						if plugin.Path != "" {
							fileName := filepath.Base(plugin.Path)
							// Remove file extension
							if idx := strings.LastIndex(fileName, "."); idx != -1 {
								fileName = fileName[:idx]
							}
							// Remove common prefixes like "lib"
							if strings.HasPrefix(fileName, "lib") {
								fileName = fileName[3:]
							}
							// Convert to uppercase and replace hyphens with underscores
							pluginName = strings.ToUpper(strings.ReplaceAll(fileName, "-", "_"))
						}
						envVarName = fmt.Sprintf("%s_%s", pluginName, strings.ToUpper(configItem.Name))
					} else {
						envVarName = fmt.Sprintf("%s_%s", strings.ToUpper(pluginType), strings.ToUpper(configItem.Name))
					}
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

// LoadGeneralPluginConfig loads configuration for a specific general plugin by finding it in the YAML config
// and setting environment variables with the extracted plugin name
func LoadGeneralPluginConfig(pluginName string) error {
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

	// Find the general plugin that matches our plugin name and set its environment variables
	if config.Plugins != nil {
		for _, plugin := range config.Plugins {
			if strings.EqualFold(plugin.Type, "general") && plugin.Path != "" {
				// Extract plugin name from path like we do in the regular function
				fileName := filepath.Base(plugin.Path)
				// Remove file extension
				if idx := strings.LastIndex(fileName, "."); idx != -1 {
					fileName = fileName[:idx]
				}
				// Remove common prefixes like "lib"
				if strings.HasPrefix(fileName, "lib") {
					fileName = fileName[3:]
				}
				// Convert to uppercase and replace hyphens with underscores
				extractedPluginName := strings.ToUpper(strings.ReplaceAll(fileName, "-", "_"))
				
				// Check if this is the plugin we're looking for
				if extractedPluginName == strings.ToUpper(pluginName) {
					for _, configItem := range plugin.Config {
						envVarName := fmt.Sprintf("%s_%s", extractedPluginName, strings.ToUpper(configItem.Name))
						// Only set if not already set (allow env vars to override config file)
						if os.Getenv(envVarName) == "" {
							os.Setenv(envVarName, configItem.Value)
						}
					}
					break // Found our plugin, no need to continue
				}
			}
		}
	}

	return nil
}

// LoadAllGeneralPluginConfigs loads configuration for all general plugins in the YAML config
// This is useful when a Go general plugin doesn't know its specific name but wants to load
// environment variables for all general plugins
func LoadAllGeneralPluginConfigs() error {
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

	// Load config for all general plugins
	if config.Plugins != nil {
		for _, plugin := range config.Plugins {
			if strings.EqualFold(plugin.Type, "general") && plugin.Path != "" {
				// Extract plugin name from path
				fileName := filepath.Base(plugin.Path)
				// Remove file extension
				if idx := strings.LastIndex(fileName, "."); idx != -1 {
					fileName = fileName[:idx]
				}
				// Remove common prefixes like "lib"
				if strings.HasPrefix(fileName, "lib") {
					fileName = fileName[3:]
				}
				// Convert to uppercase and replace hyphens with underscores
				extractedPluginName := strings.ToUpper(strings.ReplaceAll(fileName, "-", "_"))
				
				// Set environment variables for this general plugin
				for _, configItem := range plugin.Config {
					envVarName := fmt.Sprintf("%s_%s", extractedPluginName, strings.ToUpper(configItem.Name))
					// Only set if not already set (allow env vars to override config file)
					if os.Getenv(envVarName) == "" {
						os.Setenv(envVarName, configItem.Value)
					}
				}
			}
		}
	}

	return nil
}