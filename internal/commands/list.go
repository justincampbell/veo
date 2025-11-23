package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/justincampbell/veo/internal/api"
	"github.com/spf13/cobra"
)

// NewListCmd creates the list command
func NewListCmd() *cobra.Command {
	var clubSlug string
	var page int
	var all bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recordings",
		Long:  `List all recordings/matches from your Veo camera.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get auth token from environment
			token := os.Getenv("VEO_TOKEN")
			if token == "" {
				return fmt.Errorf("VEO_TOKEN environment variable is required")
			}

			// Get club slug from flag or environment variable
			if clubSlug == "" {
				clubSlug = os.Getenv("VEO_CLUB")
			}
			if clubSlug == "" {
				return fmt.Errorf("--club flag or VEO_CLUB environment variable is required")
			}

			// Create API client
			client := api.NewClient(api.WithAuthToken(token))

			// List recordings with pagination options
			opts := &api.ListRecordingsOptions{
				Page:     page,
				FetchAll: all,
			}

			recordings, err := client.ListRecordings(clubSlug, opts)
			if err != nil {
				return fmt.Errorf("failed to list recordings: %w", err)
			}

			// Print results in table format
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			fmt.Fprintln(w, "TITLE\tSLUG\tDURATION\tCREATED")
			for _, r := range recordings {
				duration := formatDuration(r.Duration)
				created := r.Created.Format("2006-01-02 15:04")
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\n", r.Title, r.Slug, duration, created)
			}
			w.Flush()

			fmt.Fprintf(os.Stderr, "\nTotal: %d recordings\n", len(recordings))

			return nil
		},
	}

	cmd.Flags().StringVarP(&clubSlug, "club", "c", "", "Club slug (or set VEO_CLUB environment variable)")
	cmd.Flags().IntVarP(&page, "page", "p", 1, "Page number (default: 1)")
	cmd.Flags().BoolVarP(&all, "all", "a", false, "Fetch all pages")

	return cmd
}

// formatDuration formats seconds into HH:MM:SS
func formatDuration(seconds int) string {
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
