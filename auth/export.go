package auth

import (
	"log"
)

// ExportPlugin is a convenience function that registers a plugin and keeps the program running
// This should be called from the plugin's main function
func ExportPlugin(plugin AuthPlugin) {
	if plugin == nil {
		log.Fatal("Cannot export nil plugin")
	}

	RegisterPlugin(plugin)

	// Validate that the plugin implements all required methods
	name := plugin.ProviderName()
	if name == "" {
		log.Fatal("Plugin must provide a non-empty provider name")
	}

	// Test basic plugin functionality
	if !plugin.HealthCheck() {
		log.Printf("Warning: Plugin %s failed health check during registration", name)
	}

	log.Printf("Go auth plugin registered: %s", name)

	// Keep the program running (required for shared library)
	select {}
}
