package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"pinboard-popular-feed/data"
)

const pinboardPopularUrl = "https://pinboard.in/popular"

func ScrapePinboardPopular() ([]*data.Bookmark, error) {
	bookmarks := make([]*data.Bookmark, 0)

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".bookmark", func(e *colly.HTMLElement) {
		id := e.Attr("id")
		title := e.ChildText(":first-child .bookmark_title")
		href := e.ChildAttr(":first-child .bookmark_title", "href")
		//fmt.Println(id)
		//fmt.Println(title)
		//fmt.Println(href)
		bookmarks = append(bookmarks, &data.Bookmark{
			Id:    id,
			Title: title,
			Url:   href,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	err := c.Visit(pinboardPopularUrl)

	return bookmarks, err
}
