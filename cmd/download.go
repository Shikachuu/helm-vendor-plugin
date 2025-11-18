package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Shikachuu/helm-vendor-plugin/internal/config"
	"github.com/Shikachuu/helm-vendor-plugin/internal/helm"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// NewDownloadCommand creates and returns a new cobra command for downloading helm charts.
// It reads the vendor charts configuration file, parses it, and downloads each specified
// helm chart to its designated destination directory.
func NewDownloadCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "download",
		Short: "Download, downloads the helm charts defined in the config file.",
		Long:  "Download, downloads the helm charts defined in the config file to their given locations.",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg, err := os.ReadFile(configPath)
			if err != nil {
				return fmt.Errorf("failed to read config file: %w", err)
			}

			cfg, err = yaml.YAMLToJSON(cfg)
			if err != nil {
				return fmt.Errorf("failed to convert yaml configuration to json: %w", err)
			}

			jcp, err := config.NewJSONConfigParser()
			if err != nil {
				return fmt.Errorf("failed to initiate json config parser: %w", err)
			}

			vcs, err := jcp.Unmarshall(cfg)
			if err != nil {
				return err
			}

			for i := range vcs {
				vc := &vcs[i]
				slog.Info("downloading chart", "repo", vc.Repository, "name", vc.Name, "destination", vc.Destination)

				fErr := helm.FetchCharts(helmCLI, vc)
				if fErr != nil {
					return fErr
				}
			}

			slog.Info("downloaded all charts", "total", len(vcs))

			return nil
		},
	}
}
