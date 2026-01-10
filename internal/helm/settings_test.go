package helm_test

import (
	"testing"

	"github.com/Shikachuu/helm-vendor-plugin/internal/helm"
	"github.com/stretchr/testify/assert"
)

func TestNewSettings(t *testing.T) {
	envs := map[string]string{
		"HELM_REGISTRY_CONFIG":   "",
		"HELM_REPOSITORY_CONFIG": "",
		"HELM_REPOSITORY_CACHE":  "",
		"HELM_CONTENT_CACHE":     "",
	}
	for k, v := range envs {
		t.Setenv(k, v)
	}

	s := helm.NewSettings()

	assert.Equal(t, s.RegistryConfig, envs["HELM_REGISTRY_CONFIG"])
	assert.Equal(t, s.RepositoryConfig, envs["HELM_REPOSITORY_CONFIG"])
	assert.Equal(t, s.RepositoryCache, envs["HELM_REPOSITORY_CACHE"])
	assert.Equal(t, s.ContentCache, envs["HELM_CONTENT_CACHE"])
	assert.Equal(t, s.Debug, envs["HELM_DEBUG"] == "1")
}
