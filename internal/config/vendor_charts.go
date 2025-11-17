// Package config responsible for configuration loading and validation
package config

// VendorChart describes a chart's properties described in the configuration file.
type VendorChart struct {
	Name        string `json:"name"`
	Repository  string `json:"repository"`
	Version     string `json:"version"`
	Destination string `json:"destination"`
	Insecure    bool   `json:"insecure"`
	Verify      bool   `json:"verify"`
}
