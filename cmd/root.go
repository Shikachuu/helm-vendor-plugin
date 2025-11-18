// Package cmd contains cobra commands for the CLI interface
package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	"helm.sh/helm/v4/pkg/cli"
)

var (
	configPath string
	helmCLI    *cli.EnvSettings
)

// NewRootCommand creates and returns the root cobra command for the helm vendor plugin.
// It configures the vendor command with subcommands and validates the configuration file path.
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "vendor [command]",
		Short:   "Vendor downloads helm charts from remote repositories.",
		Long:    "Vendor downloads helm charts from OCI and Helm repositories for vendoring, either in unpacked or tgz form.",
		Version: version,
		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			if _, err := os.Stat(configPath); os.IsNotExist(err) {
				return fmt.Errorf("config file not found: %s", configPath)
			} else if err != nil {
				return fmt.Errorf("error accessing config file: %w", err)
			}

			helmCLI = cli.New()

			ll := slog.LevelInfo
			if helmCLI.Debug {
				ll = slog.LevelDebug
			}

			slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: ll})))

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&configPath, "file", "f", ".vendor-charts.yaml", "The file contains the vendor-charts config.")

	rootCmd.AddCommand(
		NewVerifyCommand(),
		NewVersionCommand(),
		NewDownloadCommand(),
	)

	return rootCmd
}
