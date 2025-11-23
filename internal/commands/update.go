package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewUpdateCmd creates the update command
func NewUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update video metadata",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("not yet implemented - waiting for HAR analysis")
		},
	}

	return cmd
}
