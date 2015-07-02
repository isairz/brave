package brave

import (
	"sync"
)

func GetAllChapters(scraper Scraper, mangaList []MangaInfo, chapterListChan chan *[]ChapterInfo) {
	go func() {
		var wg sync.WaitGroup
		concurrency := 100
		sem := make(chan bool, concurrency)
		i := 0
		sum := 0
		for _, m := range mangaList {
			sem <- true
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				defer func() { <-sem }()
				chapterList := scraper.GetChapterList(url)
				chapterListChan <- &chapterList
				//fmt.Println(i+1, len(mangaList), chapterList[0])
				sum += len(chapterList)
				i++
			}(m.Link)
		}
		wg.Wait()
		close(chapterListChan)
		//fmt.Printf("총챕터 수 : %d\n", sum)
	}()
}
