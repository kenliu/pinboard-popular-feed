package main

import (
	"fmt"
	"github.com/gocolly/colly/v2"
)

type bookmark struct {
	id    string
	title string
	url   string
}

func ScrapePinboardPopular() []*bookmark {

	bookmarks := make([]*bookmark, 0)

	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML(".bookmark", func(e *colly.HTMLElement) {
		// fmt.Println(e)
		id := e.Attr("id")
		title := e.ChildText(":first-child .bookmark_title")
		href := e.ChildAttr(":first-child .bookmark_title", "href")
		fmt.Println(id)
		fmt.Println(title)
		fmt.Println(href)
		bookmarks = append(bookmarks, &bookmark{
			id:    id,
			title: title,
			url:   href,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("https://pinboard.in/popular/")

	return bookmarks
}
