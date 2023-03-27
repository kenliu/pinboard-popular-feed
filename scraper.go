package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

const pinboardPopularUrl = "https://pinboard.in/popular"

func ScrapePinboardPopular() []*Bookmark {
	bookmarks := make([]*Bookmark, 0)

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".bookmark", func(e *colly.HTMLElement) {
		id := e.Attr("id")
		title := e.ChildText(":first-child .bookmark_title")
		href := e.ChildAttr(":first-child .bookmark_title", "href")
		//fmt.Println(id)
		//fmt.Println(title)
		//fmt.Println(href)
		bookmarks = append(bookmarks, &Bookmark{
			id:    id,
			title: title,
			url:   href,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit(pinboardPopularUrl)

	return bookmarks
}
