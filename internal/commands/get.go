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

	cmd := &cobra.Command{
		Use:   "get <recording-id>",
		Short: "Get details for a specific recording",
		Long:  `Get detailed information about a specific recording/match.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			recordingID := args[0]

			// Get auth token from environment
			token := os.Getenv("VEO_TOKEN")
			if token == "" {
				return fmt.Errorf("VEO_TOKEN environment variable is required")
			}

			// Create API client
			client := api.NewClient(api.WithAuthToken(token))

			// Get recording details
			details, err := client.GetRecording(recordingID)
			if err != nil {
				return fmt.Errorf("failed to get recording: %w", err)
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
			printRecordingDetails(details)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&jsonOutput, "json", "j", false, "Output as JSON")

	return cmd
}

// printRecordingDetails prints recording details in a human-readable format
func printRecordingDetails(d *api.RecordingDetails) {
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
			if score, ok := stats["score"].(map[string]interface{}); ok {
				ownScore, _ := score["own"].(float64)
				oppScore, _ := score["opponent"].(float64)
				fmt.Printf("\nScore:       %.0f - %.0f\n", ownScore, oppScore)
			}

			// Age group if available
			if ageGroup, ok := d.Info["age_group"].(string); ok && ageGroup != "" {
				fmt.Printf("Age Group:   %s\n", ageGroup)
			}
		}
	}

	fmt.Printf("\nSlug:        %s\n", d.Slug)
}
