package cache

import "errors"

var (
	// ErrNotFound is returned when a cache key is not found
	ErrNotFound = errors.New("cache key not found")

	// ErrConnectionFailed is returned when connection to cache backend fails
	ErrConnectionFailed = errors.New("failed to connect to cache backend")

	// ErrInvalidTTL is returned when an invalid TTL is provided
	ErrInvalidTTL = errors.New("invalid TTL value")

	// ErrSerializationFailed is returned when serialization fails
	ErrSerializationFailed = errors.New("failed to serialize data")

	// ErrDeserializationFailed is returned when deserialization fails
	ErrDeserializationFailed = errors.New("failed to deserialize data")

	// ErrOperationFailed is returned when a cache operation fails
	ErrOperationFailed = errors.New("cache operation failed")

	// ErrNotInitialized is returned when cache is used before initialization
	ErrNotInitialized = errors.New("cache not initialized")
)