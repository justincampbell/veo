package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List videos/recordings",
		Long:  `List all videos, recordings, or matches from your Veo camera.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: Implement after analyzing HAR file
			return fmt.Errorf("not yet implemented - waiting for HAR analysis")
		},
	}

	// Add flags as needed after HAR analysis
	cmd.Flags().IntP("limit", "l", 50, "Maximum number of videos to list")
	cmd.Flags().StringP("sort", "s", "created_at", "Sort field")

	return cmd
}
