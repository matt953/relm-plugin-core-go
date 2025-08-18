package cache

/*
#include <stdlib.h>
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"unsafe"
)

var provider CacheProvider

// SetProvider sets the cache provider implementation
func SetProvider(p CacheProvider) {
	provider = p
}

//export initialize
func initialize(configBytes *C.char, configLen C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	configStr := C.GoStringN(configBytes, configLen)
	var config map[string]string
	if err := json.Unmarshal([]byte(configStr), &config); err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to parse config: %v"}`, err))
	}

	if err := provider.Initialize(config); err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to initialize: %v"}`, err))
	}

	return C.CString(`{"success": true}`)
}

//export cache_get
func cache_get(keyBytes *C.char, keyLen C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	key := C.GoStringN(keyBytes, keyLen)
	value, err := provider.Get(key)
	if err != nil {
		return C.CString("null")
	}

	// Check for cache miss (empty value)
	if value == "" {
		return C.CString("null")
	}

	// Return the raw value if it's already JSON, otherwise encode it
	if json.Valid([]byte(value)) {
		return C.CString(value)
	}
	
	jsonValue, _ := json.Marshal(value)
	return C.CString(string(jsonValue))
}

//export cache_set
func cache_set(inputBytes *C.char, inputLen C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	input := C.GoStringN(inputBytes, inputLen)
	req, err := ParseCacheSetRequest(input)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to parse request: %v"}`, err))
	}

	ttl := SecondsToDuration(req.TTL)
	if err := provider.Set(req.Key, req.Value, ttl); err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to set cache: %v"}`, err))
	}

	return C.CString(`{"success": true}`)
}

//export cache_delete
func cache_delete(keyBytes *C.char, keyLen C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	key := C.GoStringN(keyBytes, keyLen)
	deleted, err := provider.Delete(key)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to delete: %v"}`, err))
	}

	return C.CString(fmt.Sprintf("%v", deleted))
}

//export cache_delete_pattern
func cache_delete_pattern(patternBytes *C.char, patternLen C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	pattern := C.GoStringN(patternBytes, patternLen)
	count, err := provider.DeletePattern(pattern)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to delete pattern: %v"}`, err))
	}

	return C.CString(fmt.Sprintf("%d", count))
}

//export cache_exists
func cache_exists(keyBytes *C.char, keyLen C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	key := C.GoStringN(keyBytes, keyLen)
	exists, err := provider.Exists(key)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to check existence: %v"}`, err))
	}

	return C.CString(fmt.Sprintf("%v", exists))
}

//export cache_set_multiple
func cache_set_multiple(inputBytes *C.char, inputLen C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	input := C.GoStringN(inputBytes, inputLen)
	req, err := ParseCacheSetMultipleRequest(input)
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to parse request: %v"}`, err))
	}

	entries := make(map[string]string)
	for _, kv := range req.Entries {
		entries[kv.Key] = kv.Value
	}

	ttl := SecondsToDuration(req.TTL)
	if err := provider.SetMultiple(entries, ttl); err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to set multiple: %v"}`, err))
	}

	return C.CString(`{"success": true}`)
}

//export cache_stats
func cache_stats(unused1 *C.char, unused2 C.int) *C.char {
	if provider == nil {
		return C.CString(`{"error": "Cache provider not set"}`)
	}

	stats, err := provider.Stats()
	if err != nil {
		return C.CString(fmt.Sprintf(`{"error": "Failed to get stats: %v"}`, err))
	}

	jsonStats, _ := json.Marshal(stats)
	return C.CString(string(jsonStats))
}

//export provider_name
func provider_name() *C.char {
	if provider == nil {
		return C.CString("Unknown Cache Provider")
	}
	return C.CString(provider.Name())
}

//export free_result
func free_result(ptr *C.char) {
	C.free(unsafe.Pointer(ptr))
}