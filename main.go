package main

import (
	"context"
	"os"

	"github.com/Shikachuu/helm-vendor-plugin/cmd"
)

func main() {
	rootCmd := cmd.NewRootCommand()
	ctx := context.Background()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		os.Exit(1)
	}
}
