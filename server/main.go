package main

import (
	"../"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello world!")
	})

	api := router.Group("api")
	{
		api.GET("/api/MangaList", func(c *gin.Context) {
			mangaList := []brave.MangaInfo{}
			db.Find(&mangaList)
			c.String(http.StatusOK, fmt.Sprintf("%+v", mangaList))
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
			if c.Query("forced") == "true" {
				chapterList = brave.GetAllChapterList(db)
			} else {
				chapterList = brave.GetUnscrapedChapterList(db)
			}
			msg := brave.ScrapChapters(db, chapterList)
			latency := time.Since(t)
			c.String(http.StatusOK, fmt.Sprintf("%s in %s", msg, latency))
		})
	}

	router.Run(":3000")
}
