package brave

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/robertkrimen/otto"
	//"golang.org/x/net/html"

	//"fmt"
)

type Marumaru struct {
	client *http.Client
}

func NewMarumaru() *Marumaru {
	var scraper Marumaru
	scraper.initCookie()
	return &scraper
}

func (scraper *Marumaru) GetMangaList() []MangaInfo {

	doc, err := goquery.NewDocument("http://marumaru.in/c/1")
	if err != nil {
		log.Fatal(err)
	}

	list := doc.Find("#widget_bbs_review01 li")
	mangaList := make([]MangaInfo, list.Length())
	list.Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		thumbnail, _ := s.Find("img").Attr("src")
		link, _ := s.Find("a").Attr("href")

		mangaList[i] = MangaInfo{
			Name:      name,
			Thumbnail: thumbnail,
			Link:      "http://marumaru.in" + link,
		}
	})
	return mangaList
}

func (scraper *Marumaru) GetChapterList(url string) []ChapterInfo {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	content := doc.Find("#vContent")
	content.ChildrenFiltered(".snsbox").Remove()
	content.Children().Last().Remove()

	// image, _ := content.Find("img").First().Attr("src")

	list := content.Find("a")
	chapterList := make([]ChapterInfo, list.Length())
	w := 0
	list.Each(func(i int, s *goquery.Selection) {
		name := strings.TrimSpace(s.Text())
		if len(name) == 0 {
			return
		}
		link, _ := s.Attr("href")
		link = strings.Replace(link, "http://mangaumaru.com/", "http://www.mangaumaru.com/", 1)
		if !strings.HasPrefix(link, "http://www.mangaumaru.com/") {
			return
		}

		chapterList[w] = ChapterInfo{
			Name: name,
			Link: link,
		}
		w++
	})

	return chapterList[0:w]
}

func (scraper *Marumaru) GetPages(url string) []PageInfo {
	url = strings.Replace(url, "http://mangaumaru.com/", "http://www.mangaumaru.com/", 1)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36")

	res, err := scraper.client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromResponse(res)

	//pageInfo.title = strings.TrimSpace(doc.Find("article header .entry-title").Text())

	attr := "data-lazy-src"
	list := doc.Find("#content img[" + attr + "]")
	if list.Length() == 0 {
		attr = "src"
		list = doc.Find("#content img[" + attr + "]")
	}

	pageInfo := make([]PageInfo, list.Length())
	list.Each(func(i int, s *goquery.Selection) {
		src, _ := s.Attr(attr)
		pageInfo[i] = PageInfo{
			Origin: src,
		}
	})
	return pageInfo
}

func (scraper *Marumaru) initCookie() error {
	rawUrl := "http://www.mangaumaru.com/archives/189150"
	jar, _ := cookiejar.New(nil)
	scraper.client = &http.Client{
		Jar: jar,
	}

	req, err := http.NewRequest("GET", rawUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36")

	res, err := scraper.client.Do(req)
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	res.Body.Close()

	re := regexp.MustCompile("<script>(.*)</script>")
	script := re.FindStringSubmatch(string(body))[1]

	vm := otto.New()

	script = "document = {}; location = {reload: function() {}};" + script + "; document.cookie"

	value, err := vm.Run(script)
	if err != nil {
		log.Fatal(err)
	}
	rawCookie, err := value.ToString()
	if err != nil {
		log.Fatal(err)
	}
	cookie, err := makeCookie(rawCookie)
	if err != nil {
		return err
	}
	parsedUrl, err := url.Parse("http://www.mangaumaru.com/")
	if err != nil {
		return err
	}
	scraper.client.Jar.SetCookies(parsedUrl, cookie)

	return err
}
