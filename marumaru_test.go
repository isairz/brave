package brave

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMangaList(t *testing.T) {
	t.Skip()
	var scraper Scraper
	scraper = NewMarumaru()
	mangaList := scraper.GetMangaList()
	assert.NotEmpty(t, mangaList)
}

func TestChapterList(t *testing.T) {
	t.Skip()
	var scraper Scraper
	scraper = NewMarumaru()
	chapterList := scraper.GetChapterList("http://marumaru.in/?c=1/40&sort=subject&uid=45251")
	assert.NotEmpty(t, chapterList)
	fmt.Println(chapterList)
}

func TestMangaListAndChapter(t *testing.T) {
	t.Skip()
	var scraper Scraper
	scraper = NewMarumaru()

	mangaList := scraper.GetMangaList()
	chapterListChan := make(chan *[]ChapterInfo)

	go GetAllChapters(scraper, mangaList, chapterListChan)
	for range mangaList {
		<-chapterListChan
	}
}

func TestManga(t *testing.T) {
	var scraper Scraper
	scraper = NewMarumaru()
	pageInfo := scraper.GetPages("http://www.mangaumaru.com/archives/429182")
	assert.NotEmpty(t, pageInfo)
	fmt.Println(pageInfo)
}

func TestChapterListAndManga(t *testing.T) {
	t.Skip()
	// FIXME: Move to benchmark
	var scraper Scraper
	scraper = NewMarumaru()
	//chapterList := scraper.GetChapterList("http://marumaru.in/b/manga/85271")
	chapterList := scraper.GetChapterList("http://marumaru.in/b/manga/7")

	var wg sync.WaitGroup
	concurrency := 100
	sem := make(chan bool, concurrency)
	i := 0
	sum := 0
	for _, m := range chapterList {
		sem <- true
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			defer func() { <-sem }()
			pageInfo := scraper.GetPages(url)
			assert.NotEmpty(t, pageInfo)
			fmt.Println(i, len(chapterList))
			sum += len(pageInfo)
			i++
		}(m.Link)
	}
	wg.Wait()
	fmt.Printf("총페이지 수 : %d\n", sum)
}
