package main

import (
	"fmt"
	"os"

	"github.com/justincampbell/veo/internal/commands"
	"github.com/spf13/cobra"
)

var version = "dev"

func main() {
	rootCmd := &cobra.Command{
		Use:     "veo",
		Short:   "CLI for Veo sports camera",
		Long:    `A command-line interface for interacting with the Veo sports camera API.`,
		Version: version,
	}

	// Add subcommands
	rootCmd.AddCommand(commands.NewListCmd())
	rootCmd.AddCommand(commands.NewGetCmd())
	rootCmd.AddCommand(commands.NewUpdateCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
