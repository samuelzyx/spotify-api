// models/track.go
package models

import "gorm.io/gorm"

// Track represents a music track
type Track struct {
	gorm.Model
	ISRC       string
	ImageURI   string
	Title      string
	Popularity int
	ArtistID   uint   // Foreign key referencing Artist's ID
	Artist     Artist // Relationship with Artist (parent)
}
