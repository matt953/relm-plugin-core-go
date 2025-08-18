package cache

import (
	"encoding/json"
	"time"
)

// CacheProvider defines the interface for cache implementations
type CacheProvider interface {
	// Initialize the cache provider with configuration
	Initialize(config map[string]string) error

	// Get retrieves a value from the cache
	Get(key string) (string, error)

	// Set stores a value in the cache with optional TTL
	Set(key string, value string, ttl *time.Duration) error

	// Delete removes a key from the cache
	Delete(key string) (bool, error)

	// DeletePattern removes all keys matching a pattern
	DeletePattern(pattern string) (int, error)

	// Exists checks if a key exists in the cache
	Exists(key string) (bool, error)

	// SetMultiple stores multiple key-value pairs with optional TTL
	SetMultiple(entries map[string]string, ttl *time.Duration) error

	// Stats returns cache statistics
	Stats() (map[string]string, error)

	// Name returns the name of the cache provider
	Name() string
}

// CacheEntry represents a cached item with metadata
type CacheEntry struct {
	Data      string     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// CacheSetRequest represents a cache set operation
type CacheSetRequest struct {
	Key   string  `json:"key"`
	Value string  `json:"value"`
	TTL   *uint64 `json:"ttl,omitempty"`
}

// CacheSetMultipleRequest represents a batch cache set operation
type CacheSetMultipleRequest struct {
	Entries []CacheKeyValue `json:"entries"`
	TTL     *uint64         `json:"ttl,omitempty"`
}

// CacheKeyValue represents a key-value pair
type CacheKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Helper function to convert TTL seconds to duration
func SecondsToDuration(seconds *uint64) *time.Duration {
	if seconds == nil {
		return nil
	}
	duration := time.Duration(*seconds) * time.Second
	return &duration
}

// Helper function to convert duration to TTL seconds
func DurationToSeconds(duration *time.Duration) *uint64 {
	if duration == nil {
		return nil
	}
	seconds := uint64(duration.Seconds())
	return &seconds
}

// ParseCacheSetRequest parses JSON input for cache set operation
func ParseCacheSetRequest(input string) (*CacheSetRequest, error) {
	var req CacheSetRequest
	err := json.Unmarshal([]byte(input), &req)
	return &req, err
}

// ParseCacheSetMultipleRequest parses JSON input for batch cache set operation
func ParseCacheSetMultipleRequest(input string) (*CacheSetMultipleRequest, error) {
	var req CacheSetMultipleRequest
	err := json.Unmarshal([]byte(input), &req)
	return &req, err
}