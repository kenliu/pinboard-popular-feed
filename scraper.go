package main

import (
	"log/slog"
	"pinboard-popular-feed/data"

	"github.com/gocolly/colly/v2"
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
			BookmarkId: id,
			Title:      title,
			Url:        href,
		})
	})

	c.OnRequest(func(r *colly.Request) {
		slog.Debug("visiting", "url", r.URL.String())
	})

	err := c.Visit(pinboardPopularUrl)
	if err != nil {
		slog.Error("error fetching pinboard popular page", "error", err)
		return nil, err
	}

	return bookmarks, err
}
