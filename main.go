package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"pinboard-popular-feed/data"
)

func main() {
	// handle dry-run command line flag
	dryRun := flag.Bool("dryrun", false, "scan for new posts but don't post to mastodon")
	flag.Parse()

	mastodonCredentials, err := buildMastodonCredentials()
	if err != nil {
		os.Exit(1)
	}

	popular, err := fetchCurrentPinboardPopular()
	if err != nil {
		fmt.Println("error scraping popular bookmarks")
		os.Exit(1)
	}

	// initialize the bookmark store
	db := data.Init()
	db.InitStore(data.DBConfig{})

	postCount, err := postNewLinks(popular, db, mastodonCredentials, *dryRun)
	if err != nil {
		fmt.Println("error posting new links")
		os.Exit(1)
	}
	fmt.Println("posted " + fmt.Sprint(postCount) + " new bookmarks")
}

func buildMastodonCredentials() (MastodonCredentials, error) {
	if os.Getenv("MASTODON_ACCESS_TOKEN") == "" {
		fmt.Println("MASTODON_ACCESS_TOKEN not set")
		return MastodonCredentials{}, errors.New("MASTODON_ACCESS_TOKEN not set")
	}

	if os.Getenv("MASTODON_SERVER_DOMAIN") == "" {
		fmt.Println("MASTODON_SERVER_DOMAIN not set")
		return MastodonCredentials{}, errors.New("MASTODON_SERVER_DOMAIN not set")
	}

	return MastodonCredentials{
		accessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
		serverDomain: os.Getenv("MASTODON_SERVER_DOMAIN"),
	}, nil
}

func postNewLinks(popular []*data.Bookmark, db data.BookmarkStore, mastodonCredentials MastodonCredentials, dryRun bool) (int, error) {
	fmt.Println("dryrun: not posting any new bookmarks")
	var postCount int
	for i := 0; i < len(popular); i++ {
		found, err := db.FindBookmark(popular[i].Id)
		if err != nil {
			fmt.Println("error finding bookmark in store: " + popular[i].Id)
			return postCount, err
		}
		if !found {
			if dryRun {
				fmt.Println("dry run: new bookmark found, but not posted")
				continue
			} else {
				db.StoreBookmark(*popular[i])
				fmt.Println("new bookmark stored")
				postCount++
				TootBookmark(*popular[i], mastodonCredentials)
			}
		}
		// TODO implement some kind of rate limiting?
		// https://docs.joinmastodon.org/api/rate-limits/
	}
	return postCount, nil
}

func fetchCurrentPinboardPopular() ([]*data.Bookmark, error) {
	popular, _ := ScrapePinboardPopular()
	fmt.Println("current popular bookmarks: ")
	for i := 0; i < len(popular); i++ {
		fmt.Println(popular[i].Id)
		fmt.Println(popular[i].Title)
		fmt.Println(popular[i].Url)
		fmt.Println()
	}

	fmt.Println("found " + fmt.Sprint(len(popular)) + " bookmarks on pinboard popular")
	return popular, nil //TODO handle any errors
}
