package main

import (
	"fmt"
	"os"
)

type Bookmark struct {
	id    string
	title string
	url   string
}

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
		fmt.Println(popular[i].id)
		fmt.Println(popular[i].title)
		fmt.Println(popular[i].url)
		println()
	}

	TootBookmark(*popular[0], buildMastodonCredentials())
}
