package models

import "gorm.io/gorm"

type Track struct {
	gorm.Model
	ISRC       string
	ImageURI   string
	Title      string
	ArtistName string
	Popularity int
}
