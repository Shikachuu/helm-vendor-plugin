package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/Shikachuu/helm-vendor-plugin/internal/config"
	"github.com/spf13/cobra"
	"sigs.k8s.io/yaml"
)

// NewVerifyCommand creates and returns a cobra command that verifies the vendor-charts
// configuration file by reading it, converting it to JSON, and validating it against
// the expected schema.
func NewVerifyCommand() *cobra.Command {
	verifyCmd := &cobra.Command{
		Use:   "verify",
		Short: "Verifies the given vendor-charts configuration file.",
		Long:  "Verifies the given vendor-charts configuration file.",
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

			err = jcp.Validate(cfg)
			if err != nil {
				return err
			}

			slog.Info("config is valid", "path", configPath)

			return nil
		},
	}

	return verifyCmd
}
