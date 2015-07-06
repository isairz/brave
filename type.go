package brave

import (
	"github.com/jinzhu/gorm"
)

type MangaInfo struct {
	gorm.Model
	Thumbnail string
	Name      string
	Genres    []Genre
	Authors   []Author
	Artists   []Artist
	Status    string
	Type      string
	Link      string
}

type ChapterInfo struct {
	gorm.Model
	Thumbnail string
	Name      string
	Link      string
}

type PageInfo struct {
	gorm.Model
	Origin string
}

type Genre struct {
	gorm.Model
	Name string
}

type Author struct {
	gorm.Model
	Name string
}

type Artist struct {
	gorm.Model
	Name string
}
