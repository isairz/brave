package main

import (
	"../"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/itsjamie/gin-cors"
	"github.com/jinzhu/gorm"
)

func initDatabase(db *gorm.DB) {
	var err error
	*db, err = gorm.Open("mysql", "brave:brave@/brave?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
		log.Fatal(err)
	}

	brave.MigrateDatabase(db)
	db.LogMode(true)
}

func main() {
	db := &gorm.DB{}
	initDatabase(db)

	router := gin.New()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.GET("/image-proxy", func(c *gin.Context) {
		src := c.Query("src")
		res, err := brave.Proxy(src)
		if err != nil {
			// FIXME: X-box image
			c.String(404, "NotFind")
			return
		}
		w := c.Writer
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
	})

	api := router.Group("api")
	{
		api.GET("/MangaList", func(c *gin.Context) {
			q := c.Query("q")
			mangaList := []brave.MangaInfo{}
			db.Where("Name LIKE ?", "%"+q+"%").Order("Name").Limit("25").Find(&mangaList)

			result := make([]interface{}, len(mangaList))
			for i, manga := range mangaList {
				result[i] = map[string]interface{}{
					"ID":   manga.ID,
					"Name": manga.Name,
				}
			}
			c.JSON(http.StatusOK, result)
		})
		api.GET("/MangaInfo/:id", func(c *gin.Context) {
			mangaID := c.Param("id")
			mangaInfo := brave.MangaInfo{}
			db.First(&mangaInfo, mangaID)

			result := map[string]interface{}{
				"ID":   mangaInfo.ID,
				"Name": mangaInfo.Name,
			}
			c.JSON(http.StatusOK, result)
		})
		api.GET("/ChapterList/:id", func(c *gin.Context) {
			mangaID := c.Param("id")
			chapterList := []brave.ChapterInfo{}
			db.Where("manga_id=?", mangaID).Order("number").Find(&chapterList)

			result := make([]interface{}, len(chapterList))
			for i, chapter := range chapterList {
				result[i] = map[string]interface{}{
					"ID":   chapter.ID,
					"Name": chapter.Name,
				}
			}
			c.JSON(http.StatusOK, result)
		})
		api.GET("/ChapterInfo/:mangaID/:chapter", func(c *gin.Context) {
			mangaID := c.Param("mangaID")
			chapter := c.Param("chapter")
			chapterInfo := brave.ChapterInfo{}
			db.Where("manga_id=? AND number=?", mangaID, chapter).First(&chapterInfo)

			result := map[string]interface{}{
				"ID":   chapterInfo.ID,
				"Name": chapterInfo.Name,
			}
			c.JSON(http.StatusOK, result)
		})
		api.GET("/PageList/:mangaID/:chapter", func(c *gin.Context) {
			mangaID := c.Param("mangaID")
			chapter := c.Param("chapter")
			pageList := []brave.PageInfo{}
			db.Joins("JOIN chapter_infos ON page_infos.chapter_id=chapter_infos.id").Where("chapter_infos.manga_id=? AND chapter_infos.number=?", mangaID, chapter).Order("page_infos.number").Find(&pageList)

			result := make([]interface{}, len(pageList))
			for i, page := range pageList {
				result[i] = map[string]interface{}{
					"Src": page.Origin,
				}
			}
			c.JSON(http.StatusOK, result)
		})
	}

	cmd := router.Group("/admin/cmd")
	{
		cmd.GET("/ScrapMangaList", func(c *gin.Context) {
			t := time.Now()
			msg := brave.ScrapMangaList(db)
			latency := time.Since(t)
			c.String(http.StatusOK, fmt.Sprintf("%s in %s", msg, latency))
		})
		cmd.GET("/ScrapMangas", func(c *gin.Context) {
			t := time.Now()
			mangaList := brave.GetAllMangaList(db)
			msg := brave.ScrapMangas(db, mangaList)
			latency := time.Since(t)
			c.String(http.StatusOK, fmt.Sprintf("%s in %s", msg, latency))
		})
		cmd.GET("/ScrapChapters", func(c *gin.Context) {
			t := time.Now()
			var chapterList []brave.ChapterInfo
			if c.Query("force") == "true" {
				chapterList = brave.GetAllChapterList(db)
			} else {
				chapterList = brave.GetUnscrapedChapterList(db)
			}
			msg := brave.ScrapChapters(db, chapterList)
			latency := time.Since(t)
			c.String(http.StatusOK, fmt.Sprintf("%s in %s", msg, latency))
		})
	}

	router.Run(":3643")
}
