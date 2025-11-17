package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Version information variables that are set at build time via ldflags.
var (
	// version is the semantic version of the application.
	version = "v0.0.0"
	// commit is the git commit hash of the build.
	commit = ""
	// date is the build timestamp.
	date = "2025-11-17T22:05:20"
)

// NewVersionCommand creates and returns a new cobra command that prints version information.
// The version information includes the semantic version, git commit hash, and build date.
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version number",
		RunE: func(_ *cobra.Command, _ []string) error {
			fmt.Printf("Version: %s\tCommit: %s\tBuild date:%s\n", version, commit, date)
			return nil
		},
	}
}
