package main

import (
	"errors"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
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
		//log.Println(id)
		//log.Println(title)
		//log.Println(href)
		bookmarks = append(bookmarks, &data.Bookmark{
			Id:    id,
			Title: title,
			Url:   href,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL)
	})

	err := c.Visit(pinboardPopularUrl)
	if err != nil {
		log.Println("error fetching pinboard popular page: " + fmt.Sprint(errors.Unwrap(err)))
		return nil, err
	}

	return bookmarks, err
}
