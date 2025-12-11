// Package helm wraps the necessary helm interfaces to download charts
package helm

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"strings"

	"github.com/Shikachuu/helm-vendor-plugin/internal/config"
	"golang.org/x/sync/errgroup"
	"helm.sh/helm/v4/pkg/downloader"
	"helm.sh/helm/v4/pkg/getter"
	"helm.sh/helm/v4/pkg/registry"
	repo "helm.sh/helm/v4/pkg/repo/v1"
)

// FetchCharts downloads a list of VendorChart to it's location
// from it's Helm Repository or OCI Registry,
// it uses the system's repository cache and configuration.
// Repository authentication must be a separate step, this function
// already assumes you are authenticated to the given registry or repo.
//
// Returns an error if any.
func FetchCharts(s *Settings, vendorCharts []config.VendorChart) error {
	// Use getter.Getters() instead of getter.All() to avoid cli.EnvSettings dependency
	// This provides HTTP and OCI getters without pulling in Kubernetes client libraries
	getters := getter.Getters()

	rc, cErr := registry.NewClient(registry.ClientOptCredentialsFile(s.RegistryConfig))
	if cErr != nil {
		return fmt.Errorf("cannot create new OCI registry client: %w", cErr)
	}

	var eg errgroup.Group

	for i := range vendorCharts {
		eg.Go(func() error {
			logger := slog.With("name", vendorCharts[i].Name)
			logger.Info("downloading chart", "repo", vendorCharts[i].Repository, "destination", vendorCharts[i].Destination)

			err := os.MkdirAll(vendorCharts[i].Destination, 0o750)
			if err != nil {
				return fmt.Errorf("unable to create target directory: %w", err)
			}

			dl := downloader.ChartDownloader{
				Out:              os.Stdout,
				Getters:          getters,
				Verify:           getVerify(&vendorCharts[i]),
				RepositoryConfig: s.RepositoryConfig,
				RepositoryCache:  s.RepositoryCache,
				ContentCache:     s.ContentCache,
				RegistryClient:   rc,
			}

			url, err := getChartURL(getters, &vendorCharts[i])
			if err != nil {
				return fmt.Errorf("failed to get chart full URL: %w", err)
			}

			p, v, err := dl.DownloadToCache(url, vendorCharts[i].Version)
			if err != nil {
				return fmt.Errorf("unable to download chart: %w", err)
			}

			logger.Info("chart downloaded to cache", "url", url)

			if vendorCharts[i].Extract {
				logger.Info("extracting chart", "destination", vendorCharts[i].Destination)

				err = extractChartTgz(p, vendorCharts[i].Destination)
			} else {
				cURL := strings.Split(url, "/")
				destPath := path.Join(vendorCharts[i].Destination, cURL[len(cURL)-1])

				logger.Info("copying chart archive", "destination", destPath)

				err = copyChart(p, destPath)
			}

			if err != nil {
				return fmt.Errorf("unable to perform chart filemsystem action: %w", err)
			}

			if v != nil && v.SignedBy != nil {
				slog.Info("chart validated", "url", url, "hash", v.FileHash)
			}

			return nil
		})
	}

	if wErr := eg.Wait(); wErr != nil {
		return fmt.Errorf("unable to download charts: %w", wErr)
	}

	return nil
}

// getChartURL returns the full URL for the chart.
// For OCI repositories this is just Repository + Name,
// for Helm Repos it tries to find it in the registry index.
//
// Returns the full URL or an error if any.
func getChartURL(getters getter.Providers, vc *config.VendorChart) (string, error) {
	if registry.IsOCI(vc.Repository) {
		return vc.Repository + "/" + vc.Name, nil
	}

	url, err := repo.FindChartInRepoURL(
		vc.Repository,
		vc.Name,
		getters,
		repo.WithInsecureSkipTLSverify(vc.Insecure),
		repo.WithChartVersion(vc.Version),
	)
	if err != nil {
		return "", fmt.Errorf("unable to find chart in repository: %w", err)
	}

	return url, nil
}

// getVerify returns the verification status based on the VendorChart settings
//
// We must use the 2 extremes Always and Never cause that's what helms doing too.
func getVerify(vc *config.VendorChart) downloader.VerificationStrategy {
	if vc.Verify {
		return downloader.VerifyAlways
	}

	return downloader.VerifyNever
}
