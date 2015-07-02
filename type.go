package brave

import (
	"github.com/jinzhu/gorm"
)

type MangaInfo struct {
	gorm.Model
	Thumbnail string
	Name      string
	Genres    []string
	Authors   []string
	Artists   []string
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
