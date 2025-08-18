package cache

import (
	"C"
	"encoding/json"
	"fmt"
	"time"
)

// SetCachePlugin sets the global cache plugin provider that will handle all FFI calls
// This must be called in your plugin's init() function
func SetCachePlugin(provider CacheProvider) {
	SetProvider(provider)
}

// Helper macro-like function to export all required FFI functions for a cache plugin
// This should be called from your plugin's main package to ensure all functions are exported
func ExportCachePlugin() {
	// This function exists to document the required exports and ensure compilation
	// The actual exports are defined in export.go via the //export directives
}

// Helper function to convert Go map to C string for set_multiple operations
func convertEntriesToJSON(entries map[string]string, ttl *time.Duration) (string, error) {
	type entry struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	
	var entriesSlice []entry
	for k, v := range entries {
		entriesSlice = append(entriesSlice, entry{Key: k, Value: v})
	}
	
	request := struct {
		Entries []entry `json:"entries"`
		TTL     *uint64 `json:"ttl,omitempty"`
	}{
		Entries: entriesSlice,
	}
	
	if ttl != nil {
		seconds := uint64(ttl.Seconds())
		request.TTL = &seconds
	}
	
	jsonBytes, err := json.Marshal(request)
	return string(jsonBytes), err
}

// Helper function to parse cache set request from JSON
func parseCacheSetFromJSON(jsonStr string) (string, string, *time.Duration, error) {
	var req struct {
		Key   string  `json:"key"`
		Value string  `json:"value"`
		TTL   *uint64 `json:"ttl,omitempty"`
	}
	
	if err := json.Unmarshal([]byte(jsonStr), &req); err != nil {
		return "", "", nil, err
	}
	
	var ttl *time.Duration
	if req.TTL != nil {
		duration := time.Duration(*req.TTL) * time.Second
		ttl = &duration
	}
	
	return req.Key, req.Value, ttl, nil
}

// Helper function to convert C string to Go string safely
func cStringToGo(cStr *C.char, length C.int) string {
	if cStr == nil {
		return ""
	}
	return C.GoStringN(cStr, length)
}

// Helper function to convert Go string to C string
func goStringToC(str string) *C.char {
	return C.CString(str)
}

// Helper function to create error response
func createErrorResponse(msg string) *C.char {
	return C.CString(fmt.Sprintf(`{"error": "%s"}`, msg))
}

// Helper function to create success response
func createSuccessResponse() *C.char {
	return C.CString(`{"success": true}`)
}