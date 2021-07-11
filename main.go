package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func main() {
	e := gee.New()
	e.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
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
	g1.Use(func(c *gee.Context) {
		start := time.Now()
		c.Next()
		spent := time.Since(start).Nanoseconds()
		log.Println("spent=", spent)
	})
	g1.GET("/user", func(c *gee.Context) {
		c.String(http.StatusOK, "hello")
	})
	e.Static("/static","./asset")

	e.Run(":8080")
}
