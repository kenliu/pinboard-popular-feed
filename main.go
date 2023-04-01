package main

import (
	"errors"
	"fmt"
	"os"
	"pinboard-popular-feed/data"
)

func main() {
	mastodonCredentials, err := buildMastodonCredentials()
	if err != nil {
		os.Exit(1)
	}

	popular, err := fetchCurrentPinboardPopular()
	if err != nil {
		println("error scraping popular bookmarks")
		os.Exit(1)
	}

	// initialize the bookmark store
	db := data.Init()
	db.InitStore(data.DBConfig{})

	postCount, err := postNewLinks(popular, db, mastodonCredentials)
	if err != nil {
		println("error posting new links")
		os.Exit(1)
	}
	println("posted " + fmt.Sprint(postCount) + " new bookmarks")
}

func buildMastodonCredentials() (MastodonCredentials, error) {
	if os.Getenv("MASTODON_ACCESS_TOKEN") == "" {
		println("MASTODON_ACCESS_TOKEN not set")
		return MastodonCredentials{}, errors.New("MASTODON_ACCESS_TOKEN not set")
	}

	if os.Getenv("MASTODON_SERVER_DOMAIN") == "" {
		println("MASTODON_SERVER_DOMAIN not set")
		return MastodonCredentials{}, errors.New("MASTODON_SERVER_DOMAIN not set")
	}

	return MastodonCredentials{
		accessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
		serverDomain: os.Getenv("MASTODON_SERVER_DOMAIN"),
	}, nil
}

func postNewLinks(popular []*data.Bookmark, db data.BookmarkStore, mastodonCredentials MastodonCredentials) (int, error) {
	var postCount int
	for i := 0; i < len(popular); i++ {
		found, err := db.FindBookmark(popular[i].Id)
		if err != nil {
			panic(err) //TODO handle this error
		}
		if !found {
			db.StoreBookmark(*popular[i])
			println("new bookmark stored")
			postCount++
			TootBookmark(*popular[i], mastodonCredentials)
		}
		// TODO implement some kind of rate limiting?
		// https://docs.joinmastodon.org/api/rate-limits/
	}
	return postCount, nil
}

func fetchCurrentPinboardPopular() ([]*data.Bookmark, error) {
	popular := ScrapePinboardPopular()
	println("current popular bookmarks: ")
	for i := 0; i < len(popular); i++ {
		fmt.Println(popular[i].Id)
		fmt.Println(popular[i].Title)
		fmt.Println(popular[i].Url)
		println()
	}

	println("found " + fmt.Sprint(len(popular)) + " bookmarks on pinboard popular")
	return popular, nil //TODO handle any errors
}
