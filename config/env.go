package config

import (
	"os"
	"sync"
)

var (
	configOnce     sync.Once
	configInstance ConfigProvider
)

// EnvConfig implements ConfigProvider using environment variables.
type EnvConfig struct {
	apiKey      string
	baseURL     string
	model       string
	port        string
	databaseURL string
}

// NewConfigProvider instantiates a new EnvConfig provider (singleton).
func NewConfigProvider() ConfigProvider {
	configOnce.Do(func() {
		apiKey := os.Getenv("OPENROUTER_API_KEY")
		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}
		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			baseURL = "https://openrouter.ai/api/v1"
		}
		model := os.Getenv("MODEL")
		if model == "" {
			model = "openrouter/free"
		}
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		databaseURL := os.Getenv("DATABASE_URL")
		configInstance = &EnvConfig{
			apiKey:      apiKey,
			baseURL:     baseURL,
			model:       model,
			port:        port,
			databaseURL: databaseURL,
		}
	})
	return configInstance
}

// GetAPIKey returns the API Key.
func (c *EnvConfig) GetAPIKey() string {
	return c.apiKey
}

// GetBaseURL returns the endpoint Base URL.
func (c *EnvConfig) GetBaseURL() string {
	return c.baseURL
}

// GetModel returns the configured LLM model name.
func (c *EnvConfig) GetModel() string {
	return c.model
}

// GetPort returns the application port.
func (c *EnvConfig) GetPort() string {
	return c.port
}

// GetDatabaseURL returns the database connection string.
func (c *EnvConfig) GetDatabaseURL() string {
	return c.databaseURL
}
