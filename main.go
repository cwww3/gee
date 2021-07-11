package main

import (
	"fmt"
	"gee"
	"html/template"
	"log"
	"net/http"
	"time"
)

func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d-%02d-%02d", year, month, day)
}

func main() {
	e := gee.New()
	e.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})
	e.LoadHTMLGlob("templates/*")
	e.Static("/static", "./asset")

	e.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})
	e.GET("/custom", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom.tmpl", gee.H{
			"now":  time.Date(2021, 7, 11, 0, 0, 0, 0, time.Local),
			"name": "cw",
		})
	})

	e.GET("/hello", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
	})

	e.GET("/hello/:name", func(c *gee.Context) {
		// expect /hello?name=geektutu
		c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
	})

	e.GET("/assets/*filepath", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
	})

	e.POST("/login", func(c *gee.Context) {
		c.JSON(http.StatusOK, gee.H{
			"username": c.PostForm("username"),
			"password": c.PostForm("password"),
		})
	})

	g1 := e.Group("/v1")
	g1.Use(gee.Recover())
	g1.Use(func(c *gee.Context) {
		start := time.Now()
		c.Next()
		spent := time.Since(start).Nanoseconds()
		log.Println("spent=", spent)
	})
	g1.GET("/user", func(c *gee.Context) {
		c.String(http.StatusOK, "hello")
	})
	g1.GET("/panic", func(c *gee.Context) {
		panic("~~~")
		c.String(http.StatusOK, "hello")
	})

	e.Run(":8080")
}
