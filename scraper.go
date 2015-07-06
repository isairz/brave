package brave

import (
	"bufio"
	"fmt"
	"net/http"
	"strings"
)

type Scraper interface {
	GetMangaList() []MangaInfo
	GetChapterList(mangaInfo MangaInfo) MangaScraped
	GetPageList(chapterInfo ChapterInfo) ChapterScraped
}

func makeCookie(rawCookies string) ([]*http.Cookie, error) {
	rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", rawCookies)

	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
	if err != nil {
		return nil, err
	}
	cookie := req.Cookies()
	return cookie, err
}
