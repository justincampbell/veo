package models

import "time"

// Video represents a video/recording/match from Veo
// Structure will be updated after analyzing HAR file
type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Opponent    string    `json:"opponent,omitempty"`
	HomeAway    string    `json:"home_away,omitempty"` // "home" or "away"
	MatchType   string    `json:"match_type,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	SharingURL  string    `json:"sharing_url,omitempty"`
	HighlightsURL string  `json:"highlights_url,omitempty"`
}

// UpdateVideoRequest represents a request to update video metadata
type UpdateVideoRequest struct {
	Title     *string `json:"title,omitempty"`
	Opponent  *string `json:"opponent,omitempty"`
	HomeAway  *string `json:"home_away,omitempty"`
	MatchType *string `json:"match_type,omitempty"`
}

// Sides represents the team sides for goal tracking
// Structure will be updated after analyzing HAR file
type Sides struct {
	Home string `json:"home"`
	Away string `json:"away"`
}

// UpdateSidesRequest represents a request to update team sides
type UpdateSidesRequest struct {
	Sides Sides `json:"sides"`
}
