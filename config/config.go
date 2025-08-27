package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
)

var (
	onceConfig sync.Once
	cfg        *Config
)

// Config provides thread-safe access to a TOML configuration
type Config struct {
	data map[string]any
	path string
	mu   sync.RWMutex // Protects data for thread-safe access
}

func LoadConfig() *Config {
	onceConfig.Do(func() {
		var err error
		cfg, err = NewConfig("aur-pkg-helper.toml")
		if err != nil {
			slog.Error("Could not read config file", slog.Any("error", err))
			cfg = nil
		}
	})

	return cfg
}

// NewConfig loads a TOML file into a Config
func NewConfig(filename string) (*Config, error) {
	configPath, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("error getting home config directory: %w", err)
	}

	configFile := filepath.Clean(filepath.Join(configPath, filename))

	// Avoid directory traversal attack
	if !strings.HasPrefix(configFile, configPath) {
		return nil, errors.New("invalid config path")
	}

	// Avoid directory traversal attack
	if rel, err := filepath.Rel(configPath, configFile); err != nil || strings.HasPrefix(rel, "..") || filepath.IsAbs(rel) {
		return nil, errors.New("invalid config file path: potential directory traversal")
	}

	info, err := os.Stat(configFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist: %s", configFile)
		}

		return nil, fmt.Errorf("error checking config file: %w", err)
	}

	if !info.Mode().IsRegular() {
		return nil, fmt.Errorf("config file is not a regular file: %s", configFile)
	}

	data := map[string]any{}
	if _, err := toml.DecodeFile(configFile, &data); err != nil {
		return nil, fmt.Errorf("error loading TOML file: %w", err)
	}

	return &Config{data: data, path: configFile}, nil
}

// Reload re-parses the TOML file
func (c *Config) Reload() error {
	if c == nil {
		return errors.New("invalid configuration")
	}

	newData := map[string]any{}
	if _, err := toml.DecodeFile(c.path, &newData); err != nil {
		return fmt.Errorf("error reloading TOML file: %w", err)
	}

	c.mu.Lock()
	c.data = newData
	c.mu.Unlock()

	return nil
}

// String gets a string value with an optional default
func (c *Config) String(dottedKey, defaultVal string) string {
	if c == nil {
		slog.Error("Configuration map is nil")
		return ""
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, err := getNestedValue(c.data, dottedKey)
	if err != nil || value == nil {
		return defaultVal
	}

	str, ok := value.(string)
	if !ok {
		return defaultVal
	}
	return str
}

// Int gets an int value with an optional default
func (c *Config) Float64(dottedKey string, defaultVal float64) float64 {
	if c == nil {
		slog.Error("Configuration map is nil")
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, err := getNestedValue(c.data, dottedKey)
	if err != nil || value == nil {
		return defaultVal
	}

	// TOML numbers are parsed as float64
	f, ok := value.(float64)
	if !ok {
		return defaultVal
	}

	return f
}

func (c *Config) Int(dottedKey string, defaultVal int) int {
	if c == nil {
		slog.Error("Configuration map is nil")
		return 0
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	value, err := getNestedValue(c.data, dottedKey)
	if err != nil || value == nil {
		return defaultVal
	}

	// TOML numbers are parsed as float64
	f, ok := value.(float64)
	if !ok {
		return defaultVal
	}

	return int(f)
}

// getNestedValue retrieves a value from the nested map using dotted notation
func getNestedValue(data map[string]any, dottedKey string) (any, error) {
	parts := strings.Split(dottedKey, ".")
	current := data

	for i, part := range parts[:len(parts)-1] {
		next, ok := current[part]
		if !ok {
			return nil, fmt.Errorf("section %q not found", strings.Join(parts[:i+1], "."))
		}

		current, ok = next.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("section %q is not a map", strings.Join(parts[:i+1], "."))
		}
	}

	value, exists := current[parts[len(parts)-1]]
	if !exists {
		return nil, fmt.Errorf("key %q not found", dottedKey)
	}

	return value, nil
}
