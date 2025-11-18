// Package helm wraps the necessary helm interfaces to download charts
package helm

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Shikachuu/helm-vendor-plugin/internal/config"
	"helm.sh/helm/v4/pkg/cli"
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
func FetchCharts(s *cli.EnvSettings, repoSettings *config.VendorChart) error {
	err := os.MkdirAll(repoSettings.Destination, 0o750)
	if err != nil {
		return fmt.Errorf("unable to create target directory: %w", err)
	}

	getters := getter.All(s)

	verify := downloader.VerifyNever
	if repoSettings.Verify {
		verify = downloader.VerifyAlways
	}

	rc, err := registry.NewClient(registry.ClientOptCredentialsFile(s.RegistryConfig))
	if err != nil {
		return fmt.Errorf("cannot create new OCI registry client: %w", err)
	}

	dl := downloader.ChartDownloader{
		Out:              os.Stdout,
		Getters:          getters,
		Verify:           verify,
		RepositoryConfig: s.RepositoryConfig,
		RepositoryCache:  s.RepositoryCache,
		ContentCache:     s.ContentCache,
		RegistryClient:   rc,
	}

	url := repoSettings.Repository + "/" + repoSettings.Name

	if !registry.IsOCI(repoSettings.Repository) {
		url, err = repo.FindChartInRepoURL(
			repoSettings.Repository,
			repoSettings.Name,
			getters,
			repo.WithInsecureSkipTLSverify(repoSettings.Insecure),
			repo.WithChartVersion(repoSettings.Version),
		)
		if err != nil {
			return fmt.Errorf("unable to find chart in repository: %w", err)
		}
	}

	p, v, err := dl.DownloadTo(url, repoSettings.Version, repoSettings.Destination)
	if err != nil {
		return fmt.Errorf("unable to download chart: %w", err)
	}

	slog.Info("chart downloaded", "url", url, "destination", p)

	if v != nil {
		slog.Info("chart validated", "url", url, "hash", v.FileHash)
	}

	return nil
}
