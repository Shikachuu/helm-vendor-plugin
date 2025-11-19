// Package helm wraps the necessary helm interfaces to download charts
package helm

import "os"

// Settings holds the Helm configuration paths needed for downloading charts.
// These are read from HELM_* environment variables set by Helm when running plugins.
type Settings struct {
	RegistryConfig   string
	RepositoryConfig string
	RepositoryCache  string
	ContentCache     string
	Debug            bool
}

// NewSettings creates Settings from Helm environment variables.
// Helm automatically sets these when running a plugin.
func NewSettings() *Settings {
	s := &Settings{
		RegistryConfig:   os.Getenv("HELM_REGISTRY_CONFIG"),
		RepositoryConfig: os.Getenv("HELM_REPOSITORY_CONFIG"),
		RepositoryCache:  os.Getenv("HELM_REPOSITORY_CACHE"),
		ContentCache:     os.Getenv("HELM_CONTENT_CACHE"),
		Debug:            false,
	}

	if os.Getenv("HELM_DEBUG") == "1" {
		s.Debug = true
	}

	return s
}

