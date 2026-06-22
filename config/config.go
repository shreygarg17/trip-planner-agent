package config

// ConfigProvider defines the interface to access configuration variables.
type ConfigProvider interface {
	GetAPIKey() string
	GetBaseURL() string
	GetModel() string
	GetPort() string
	GetDatabaseURL() string
}
