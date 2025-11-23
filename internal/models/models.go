package models

import "time"

// Recording represents a recording from the list endpoint
type Recording struct {
	Camera       string    `json:"camera"`
	Created      time.Time `json:"created"`  // Upload/processing time
	Start        time.Time `json:"start"`    // Actual recording start time (match date)
	Duration     int       `json:"duration"` // in seconds
	Identifier   string    `json:"identifier"`
	Slug         string    `json:"slug"`
	Title        string    `json:"title"`
	URL          string    `json:"url"`
	Thumbnail    string    `json:"thumbnail"`
	ReelURL      string    `json:"reel_url"` // Full game highlights/reel download URL
	Team         string    `json:"team"`
	Privacy      string    `json:"privacy"`
	Permissions  string    `json:"permissions"`
	IsAccessible bool      `json:"is_accessible"`
}
