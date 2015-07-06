package brave

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	db.AutoMigrate(&MangaInfo{})
	db.AutoMigrate(&ChapterInfo{})
	db.AutoMigrate(&PageInfo{})
}

func GetMangaList(db *gorm.DB) []MangaInfo {
	var mangaList []MangaInfo
	db.Find(&mangaList)
	return mangaList
}

var scraper Scraper = NewMarumaru()

func GetAllMangaList(db *gorm.DB) (result []MangaInfo) {
	db.Find(&result)
	return
}

func GetUnscrapedChapterList(db *gorm.DB) (result []ChapterInfo) {
	db.Where("status = ''").Find(&result)
	return
}

func ScrapMangaList(db *gorm.DB) string {
	mangaList := scraper.GetMangaList()

	for _, mangaInfo := range mangaList {
		var f MangaInfo
		db.Where(&MangaInfo{Link: mangaInfo.Link}).First(&f)
		if f.ID == 0 {
			db.Create(&mangaInfo)
		} else {
			//db.Model(&f).Update(&mangaInfo)
		}
	}
	return fmt.Sprintf("%d개의 만화를 스크랩", len(mangaList))
}

func ScrapMangas(db *gorm.DB, mangaList []MangaInfo) string {
	ch := make(chan MangaScraped)
	go GetAllChapters(scraper, mangaList, ch)

	NumberOfChapter := 0
	NumberOfNewChapter := 0
	for range mangaList {
		scraped := <-ch
		NumberOfChapter += len(scraped.ChapterList)
		db.Model(&scraped.Original).Update(&scraped.Additional)
		for _, chapterInfo := range scraped.ChapterList {
			var f ChapterInfo
			db.Where(&ChapterInfo{Link: chapterInfo.Link}).First(&f)
			if f.ID == 0 {
				db.Create(&chapterInfo)
				NumberOfNewChapter++
			} else {
				db.Model(&f).Update(&chapterInfo)
			}
		}
	}

	return fmt.Sprintf("%d개의 만화에서 %d(+%d)개의 챕터", len(mangaList), NumberOfChapter, NumberOfNewChapter)
}

func ScrapChapters(db *gorm.DB, chapterList []ChapterInfo) string {
	ch := make(chan ChapterScraped)
	go GetAllPages(scraper, chapterList, ch)

	NumberOfNewChapter := len(chapterList)
	NumberOfNewPage := 0
	for range chapterList {
		scraped := <-ch
		db.Find(&scraped.Original).Update(&scraped.Additional)
		db.Where(&PageInfo{ChapterID: scraped.Original.ID}).Delete(&PageInfo{})
		for _, pageInfo := range scraped.PageList {
			db.Create(&pageInfo)
			NumberOfNewPage++
		}
	}

	return fmt.Sprintf("%d개의 챕터에서 %d개의 페이지를 스크랩", NumberOfNewChapter, NumberOfNewPage)
}
