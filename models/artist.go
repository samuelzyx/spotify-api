// models/artist.go
package models

import "gorm.io/gorm"

// Artist represents an artist in the music system
type Artist struct {
	gorm.Model
	Name   string
	Tracks []Track // One-to-many relationship with Track
}
