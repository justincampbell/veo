package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/justincampbell/veo/internal/api"
	"github.com/spf13/cobra"
)

// NewGetCmd creates the get command
func NewGetCmd() *cobra.Command {
	var jsonOutput bool
	var clubSlug string

	cmd := &cobra.Command{
		Use:   "get <recording-id|latest>",
		Short: "Get details for a specific recording",
		Long:  `Get detailed information about a specific recording/match.

Use "latest" to get the most recent recording.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			recordingID := args[0]

			// Get auth token from environment
			token := os.Getenv("VEO_TOKEN")
			if token == "" {
				return fmt.Errorf("VEO_TOKEN environment variable is required")
			}

			// Create API client
			client := api.NewClient(api.WithAuthToken(token))

			// Handle "latest" special case
			if recordingID == "latest" {
				// Get club slug from flag or environment variable
				if clubSlug == "" {
					clubSlug = os.Getenv("VEO_CLUB")
				}
				if clubSlug == "" {
					return fmt.Errorf("--club flag or VEO_CLUB environment variable is required for 'latest'")
				}

				// List recordings to get the latest one
				opts := &api.ListRecordingsOptions{Page: 1}
				result, err := client.ListRecordings(clubSlug, opts)
				if err != nil {
					return fmt.Errorf("failed to list recordings: %w", err)
				}

				if len(result.Recordings) == 0 {
					return fmt.Errorf("no recordings found")
				}

				// Use the first recording (most recent)
				recordingID = result.Recordings[0].Identifier
			}

			// Get recording details
			details, err := client.GetRecording(recordingID)
			if err != nil {
				return fmt.Errorf("failed to get recording: %w", err)
			}

			// Get periods for kickoff timestamp
			periods, err := client.GetPeriods(details.Slug)
			if err != nil {
				// Don't fail if periods aren't available, just log
				fmt.Fprintf(os.Stderr, "Warning: could not fetch periods: %v\n", err)
			}

			// Output as JSON if requested
			if jsonOutput {
				encoder := json.NewEncoder(os.Stdout)
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(details); err != nil {
					return fmt.Errorf("failed to encode JSON: %w", err)
				}
				return nil
			}

			// Print human-readable format
			printRecordingDetails(details, periods)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")
	cmd.Flags().StringVarP(&clubSlug, "club", "c", "", "Club slug (required for 'latest', or set VEO_CLUB environment variable)")

	return cmd
}

// printRecordingDetails prints recording details in a human-readable format
func printRecordingDetails(d *api.RecordingDetails, periods []api.Period) {
	fmt.Printf("ID:          %s\n", d.Identifier)
	fmt.Printf("Title:       %s\n", d.Title)
	fmt.Printf("Type:        %s\n", d.Type)

	// Convert times to local timezone
	startLocal := d.Start.Local()
	fmt.Printf("Start:       %s\n", startLocal.Format("2006-01-02 15:04:05 MST"))

	endLocal := d.End.Local()
	fmt.Printf("End:         %s\n", endLocal.Format("2006-01-02 15:04:05 MST"))

	duration := formatDuration(d.Duration)
	fmt.Printf("Duration:    %s\n", duration)

	// Team information
	if d.OwnTeamHomeOrAway != "" {
		fmt.Printf("\nTeam:        %s\n", d.OwnTeamHomeOrAway)
	}
	if d.OwnTeamColor != "" {
		fmt.Printf("Own Color:   %s\n", d.OwnTeamColor)
	}
	if d.OwnTeamFormation != "" {
		fmt.Printf("Formation:   %s\n", d.OwnTeamFormation)
	}

	// Opponent information
	if d.OpponentTeamName != "" || d.OpponentClubName != "" {
		fmt.Printf("\nOpponent:    %s", d.OpponentTeamName)
		if d.OpponentClubName != "" && d.OpponentClubName != d.OpponentTeamName {
			fmt.Printf(" (%s)", d.OpponentClubName)
		}
		fmt.Println()
	}
	if d.OpponentTeamColor != "" {
		fmt.Printf("Opp Color:   %s\n", d.OpponentTeamColor)
	}
	if d.OpponentShortName != "" {
		fmt.Printf("Opp Short:   %s\n", d.OpponentShortName)
	}
	if d.OpponentTeamFormation != "" {
		fmt.Printf("Opp Form:    %s\n", d.OpponentTeamFormation)
	}

	// Score if available
	if d.Info != nil {
		if stats, ok := d.Info["stats"].(map[string]interface{}); ok {
			// Try score_aggregated first (actual final score), fall back to score
			var ownScore, oppScore float64
			var hasScore bool

			if scoreAgg, ok := stats["score_aggregated"].(map[string]interface{}); ok {
				if own, ok := scoreAgg["own"].(float64); ok {
					ownScore = own
					hasScore = true
				}
				if opp, ok := scoreAgg["opponent"].(float64); ok {
					oppScore = opp
					hasScore = true
				}
			} else if score, ok := stats["score"].(map[string]interface{}); ok {
				if own, ok := score["own"].(float64); ok {
					ownScore = own
					hasScore = true
				}
				if opp, ok := score["opponent"].(float64); ok {
					oppScore = opp
					hasScore = true
				}
			}

			if hasScore {
				fmt.Printf("Score:       %.0f-%.0f\n", ownScore, oppScore)
			}

			// Age group if available
			if ageGroup, ok := d.Info["age_group"].(string); ok && ageGroup != "" {
				fmt.Printf("Age Group:   %s\n", ageGroup)
			}
		}
	}

	fmt.Printf("\nSlug:        %s\n", d.Slug)

	// Share URL with kickoff timestamp
	if len(periods) > 0 && len(periods[0].Timeframe) > 0 {
		kickoffSeconds := periods[0].Timeframe[0]
		kickoffTime := formatTimestamp(kickoffSeconds)
		shareURL := fmt.Sprintf("https://app.veo.co/matches/%s/#t=%s", d.Slug, kickoffTime)
		fmt.Printf("\nShare URL:   %s\n", shareURL)
	} else {
		// Fallback to share URL without timestamp
		shareURL := fmt.Sprintf("https://app.veo.co/matches/%s/", d.Slug)
		fmt.Printf("\nShare URL:   %s\n", shareURL)
	}

	// Highlights URL
	if d.ReelURL != "" {
		fmt.Printf("Highlights:  %s\n", d.ReelURL)
	}
}

// formatTimestamp converts seconds to MM:SS format for URL timestamps
func formatTimestamp(seconds int) string {
	minutes := seconds / 60
	secs := seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, secs)
}
