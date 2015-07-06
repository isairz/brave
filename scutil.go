package brave

import (
	"sync"
)

func NormalizeStatus(s string) string {
	switch s {
	case "주간":
		return "Weekly"
	case "격주":
		return "Biweekly"
	case "월간":
		return "Monthly"
	case "격월/비정기":
		return "Ongoing"
	case "완결", "붕탁 완결":
		return "Completed"
	case "단행본":
	case "단편":
		return "Completed"
	case "와이!":
	case "오토코노코 앤솔로지":
	case "여장소년 엔솔로지":
	case "오토코노코타임":
	}
	return ""
}

func GetAllChapters(scraper Scraper, mangaList []MangaInfo, ch chan MangaScraped) {
	go func() {
		var wg sync.WaitGroup
		concurrency := 100
		sem := make(chan bool, concurrency)
		for _, mangaInfo := range mangaList {
			sem <- true
			wg.Add(1)
			go func(mangaInfo MangaInfo) {
				defer wg.Done()
				defer func() { <-sem }()
				ch <- scraper.GetChapterList(mangaInfo)
			}(mangaInfo)
		}
		wg.Wait()
		close(ch)
	}()
}

func GetAllPages(scraper Scraper, chapterList []ChapterInfo, ch chan ChapterScraped) {
	go func() {
		var wg sync.WaitGroup
		concurrency := 100
		sem := make(chan bool, concurrency)
		for _, chapterInfo := range chapterList {
			sem <- true
			wg.Add(1)
			go func(chapterInfo ChapterInfo) {
				defer wg.Done()
				defer func() { <-sem }()
				ch <- scraper.GetPageList(chapterInfo)
			}(chapterInfo)
		}
		wg.Wait()
		close(ch)
	}()
}
