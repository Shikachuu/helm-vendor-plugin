package helm

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Shikachuu/helm-vendor-plugin/internal/config"
	"github.com/stretchr/testify/require"
	"helm.sh/helm/v4/pkg/downloader"
	"helm.sh/helm/v4/pkg/getter"
)

func TestGetVerify(t *testing.T) {
	tests := []struct {
		vc   *config.VendorChart
		name string
		want downloader.VerificationStrategy
	}{
		{
			name: "verify enabled",
			vc: &config.VendorChart{
				Name:       "test-chart",
				Repository: "https://example.com/charts",
				Version:    "1.0.0",
				Verify:     true,
			},
			want: downloader.VerifyAlways,
		},
		{
			name: "verify disabled",
			vc: &config.VendorChart{
				Name:       "test-chart",
				Repository: "https://example.com/charts",
				Version:    "1.0.0",
				Verify:     false,
			},
			want: downloader.VerifyNever,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getVerify(tt.vc)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestGetChartURL(t *testing.T) {
	tests := []struct {
		vc      *config.VendorChart
		name    string
		want    string
		errMsg  string
		wantErr bool
	}{
		{
			name: "OCI repository",
			vc: &config.VendorChart{
				Name:       "traefik",
				Repository: "oci://ghcr.io/traefik/helm",
				Version:    "37.4.0",
			},
			want:    "oci://ghcr.io/traefik/helm/traefik",
			wantErr: false,
		},
		{
			name: "OCI repository with registry.io",
			vc: &config.VendorChart{
				Name:       "nginx",
				Repository: "oci://registry.example.com/charts",
				Version:    "1.0.0",
			},
			want:    "oci://registry.example.com/charts/nginx",
			wantErr: false,
		},
		{
			name: "OCI repository with docker.io",
			vc: &config.VendorChart{
				Name:       "mychart",
				Repository: "oci://docker.io/myrepo",
				Version:    "2.0.0",
			},
			want:    "oci://docker.io/myrepo/mychart",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For OCI URLs, getters are not used (function returns early)
			// Pass nil to avoid any potential network calls or dependencies
			got, err := getChartURL(nil, tt.vc)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)

				return
			}

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
func TestGetChartURL_HTTPRepository(t *testing.T) {
	tests := []struct {
		name      string
		chartName string
		version   string
		indexYAML string
		wantURL   string
		errMsg    string
		wantErr   bool
	}{
		{
			name:      "successful chart lookup",
			chartName: "mychart",
			version:   "1.0.0",
			indexYAML: `apiVersion: v1
entries:
  mychart:
    - name: mychart
      version: 1.0.0
      urls:
        - charts/mychart-1.0.0.tgz`,
			wantURL: "/charts/mychart-1.0.0.tgz",
			wantErr: false,
		},
		{
			name:      "chart with absolute URL",
			chartName: "nginx",
			version:   "2.0.0",
			indexYAML: `apiVersion: v1
entries:
  nginx:
    - name: nginx
      version: 2.0.0
      urls:
        - https://example.com/charts/nginx-2.0.0.tgz`,
			wantURL: "https://example.com/charts/nginx-2.0.0.tgz",
			wantErr: false,
		},
		{
			name:      "chart not found",
			chartName: "nonexistent",
			version:   "1.0.0",
			indexYAML: `apiVersion: v1
entries:
  mychart:
    - name: mychart
      version: 1.0.0
      urls:
        - charts/mychart-1.0.0.tgz`,
			wantErr: true,
			errMsg:  "unable to find chart in repository",
		},
		{
			name:      "version not found",
			chartName: "mychart",
			version:   "2.0.0",
			indexYAML: `apiVersion: v1
entries:
  mychart:
    - name: mychart
      version: 1.0.0
      urls:
        - charts/mychart-1.0.0.tgz`,
			wantErr: true,
			errMsg:  "unable to find chart in repository",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test HTTP server that serves the index.yaml
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/index.yaml" {
					w.Header().Set("Content-Type", "application/x-yaml")
					_, _ = w.Write([]byte(tt.indexYAML))

					return
				}

				http.NotFound(w, r)
			}))
			defer server.Close()

			vc := &config.VendorChart{
				Name:       tt.chartName,
				Repository: server.URL,
				Version:    tt.version,
				Insecure:   false,
			}

			getters := getter.Getters()
			got, err := getChartURL(getters, vc)

			if tt.wantErr {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.errMsg)

				return
			}

			require.NoError(t, err)

			// For relative URLs, the test server URL will be prepended
			if tt.wantURL[0] == '/' {
				require.Equal(t, server.URL+tt.wantURL, got)
			} else {
				require.Equal(t, tt.wantURL, got)
			}
		})
	}
}
