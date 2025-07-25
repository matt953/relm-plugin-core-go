# Relm Plugin Core - Go

Go library for developing Relm plugins. This library provides the FFI bindings and interfaces needed to create Go plugins that can be loaded by the Relm Core application.

## Overview

This library enables Go developers to create plugins for Relm by providing:

- **Plugin Interfaces** - Standard interfaces for storage and general plugins
- **FFI Bindings** - C-compatible exports for integration with Rust core
- **Error Handling** - Standardized error types and handling
- **Memory Management** - Safe memory handling for FFI operations

## Plugin Types

### Storage Plugins
Storage plugins provide file storage capabilities for Relm Core. Only one storage plugin can be active at a time.

### General Plugins
General plugins provide event-driven functionality and can respond to various system events. **Multiple general plugins can be loaded simultaneously**, allowing you to compose different functionalities (e.g., notifications, analytics, logging) from separate plugins.

## Usage

```go
package main

import (
    "github.com/realm/relm-plugin-core-go/storage"
)

type MyStoragePlugin struct {
    // Your storage implementation
}

func (p *MyStoragePlugin) StoreFile(path string, data []byte, contentType *string) error {
    // Implement file storage
    return nil
}

func (p *MyStoragePlugin) RetrieveFile(path string) ([]byte, error) {
    // Implement file retrieval
    return nil, nil
}

// ... implement other StoragePlugin methods

func main() {
    plugin := &MyStoragePlugin{}
    storage.RegisterPlugin(plugin)
}
```

## Building Plugins

Plugins must be built as C shared libraries:

```bash
go build -buildmode=c-shared -o plugin.so main.go
```

## FFI Interface

The library automatically exports these C-compatible functions:

- `store_file_with_content_type`
- `store_file`
- `retrieve_file` 
- `delete_file`
- `file_exists`
- `generate_file_url`
- `provider_name`
- `init_plugin`

## Error Handling

Use the provided error types for consistent error reporting:

```go
import "github.com/realm/relm-plugin-core-go/types"

return types.NewStorageError("Failed to connect to storage backend")
```

## Memory Safety

The library handles all FFI memory management automatically. Plugin developers should focus on their storage logic without worrying about C memory management.

## Development

This library is designed to work with the Relm Plugin Development Kit (PDK). Use the PDK to scaffold new plugins:

```bash
relm-pdk new storage my-plugin --lang go
```