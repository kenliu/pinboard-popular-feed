package main

import (
	"fmt"
	"os"
	"pinboard-popular-feed/data"
)

func buildMastodonCredentials() MastodonCredentials {
	if os.Getenv("MASTODON_ACCESS_TOKEN") == "" {
		println("MASTODON_ACCESS_TOKEN not set")
		os.Exit(1)
	}

	if os.Getenv("MASTODON_SERVER_DOMAIN") == "" {
		println("MASTODON_SERVER_DOMAIN not set")
		os.Exit(1)
	}

	return MastodonCredentials{
		accessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
		serverDomain: os.Getenv("MASTODON_SERVER_DOMAIN"),
	}
}

func main() {
	popular := ScrapePinboardPopular()
	println("current popular bookmarks: ")
	for i := 0; i < len(popular); i++ {
		fmt.Println(popular[i].Id)
		fmt.Println(popular[i].Title)
		fmt.Println(popular[i].Url)
		println()
	}

	println("found " + fmt.Sprint(len(popular)) + " bookmarks on pinboard popular")

	// initialize the bookmark store
	db := data.Init()
	db.InitStore(data.DBConfig{})

	var postCount int
	for i := 0; i < len(popular); i++ {
		found, err := db.FindBookmark(popular[i].Id)
		if err != nil {
			panic(err)
		}
		if !found {
			db.StoreBookmark(*popular[i])
			println("new bookmark stored")
			postCount++
			TootBookmark(*popular[i], buildMastodonCredentials())
		}
		// TODO implement some kind of rate limiting?
		// https://docs.joinmastodon.org/api/rate-limits/
	}
	println("posted " + fmt.Sprint(postCount) + " new bookmarks")
}
