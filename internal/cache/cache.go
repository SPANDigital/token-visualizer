package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Cache provides simple disk-based caching
type Cache struct {
	dir string
	mu  sync.RWMutex
}

// NewCache creates a new cache with the specified directory
func NewCache(dir string) (*Cache, error) {
	if dir == "" {
		// Use default cache directory
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		dir = filepath.Join(homeDir, ".cache", "token-visualizer")
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory: %w", err)
	}

	return &Cache{dir: dir}, nil
}

// Get retrieves a value from cache
func (c *Cache) Get(key string, value interface{}) error {
	c.mu.RLock()
	defer c.mu.RUnlock()

	filePath := c.keyToPath(key)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("cache miss: %w", err)
		}
		return fmt.Errorf("failed to read cache file: %w", err)
	}

	if err := json.Unmarshal(data, value); err != nil {
		return fmt.Errorf("failed to unmarshal cache data: %w", err)
	}

	return nil
}

// Set stores a value in cache
func (c *Cache) Set(key string, value interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache data: %w", err)
	}

	filePath := c.keyToPath(key)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// keyToPath converts a cache key to a file path
func (c *Cache) keyToPath(key string) string {
	hash := sha256.Sum256([]byte(key))
	hashStr := hex.EncodeToString(hash[:])
	return filepath.Join(c.dir, hashStr+".json")
}

// Clear removes all cached entries
func (c *Cache) Clear() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return os.RemoveAll(c.dir)
}
