package brave

import (
	"github.com/jinzhu/gorm"
)

/* Types for Database */
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
	Status    string
	Link      string `sql:"index"; unique_index`
}

type PageInfo struct {
	gorm.Model
	MangaID   uint `sql:"index"`
	ChapterID uint `sql:"index"`
	Number    uint
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

/* Types for scrap */
type MangaScraped struct {
	Original    MangaInfo
	Additional  MangaInfo
	ChapterList []ChapterInfo
}

type ChapterScraped struct {
	Original   ChapterInfo
	Additional ChapterInfo
	PageList   []PageInfo
}
