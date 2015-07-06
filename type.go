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
	Link      string `sql:"index"; unique_index`
}

type ChapterInfo struct {
	gorm.Model
	MangaID   uint `sql:"index"`
	Thumbnail string
	Name      string
	Link      string `sql:"index"; unique_index`
}

type PageInfo struct {
	gorm.Model
	MangaID   uint `sql:"index"`
	ChapterID uint `sql:"index"`
	Origin    string
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
