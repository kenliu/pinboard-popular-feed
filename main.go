package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"pinboard-popular-feed/data"
)

func main() {
	// handle dry-run command line flag
	dryRun := flag.Bool("dryrun", false, "scan for new posts but don't post to mastodon")
	flag.Parse()

	// set up logging
	logFile, err := os.OpenFile("pinboard-popular-feed.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("error opening log file")
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("starting pinboard-popular-feed processing")
	if *dryRun {
		log.Println("DRY RUN MODE: Showing what would be posted")
	}

	// set up mastodon credentials
	mastodonCredentials, err := buildMastodonCredentials()
	if err != nil {
		log.Println("error setting up mastodon credentials:", err)
		os.Exit(1)
	}

	popular, err := fetchCurrentPinboardPopular()
	if err != nil {
		log.Println("error scraping popular bookmarks:", err)
		os.Exit(1)
	}

	// initialize the bookmark store
	var db = data.BookmarkStore{}
	db.InitStore(data.CreateDBConfigFromEnv())

	postCount, err := postNewLinks(popular, db, mastodonCredentials, *dryRun)
	if err != nil {
		log.Println("error posting new bookmarks:", err)
		os.Exit(1)
	}

	log.Println("posted", postCount, "new bookmarks")
	log.Println("finished pinboard-popular-feed processing")
}

func buildMastodonCredentials() (MastodonCredentials, error) {
	if os.Getenv("MASTODON_ACCESS_TOKEN") == "" {
		log.Println("MASTODON_ACCESS_TOKEN not set")
		return MastodonCredentials{}, errors.New("MASTODON_ACCESS_TOKEN not set")
	}

	if os.Getenv("MASTODON_SERVER_DOMAIN") == "" {
		log.Println("MASTODON_SERVER_DOMAIN not set")
		return MastodonCredentials{}, errors.New("MASTODON_SERVER_DOMAIN not set")
	}

	return MastodonCredentials{
		accessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
		serverDomain: os.Getenv("MASTODON_SERVER_DOMAIN"),
	}, nil
}

func postNewLinks(popular []*data.Bookmark, db data.BookmarkStore, mastodonCredentials MastodonCredentials, dryRun bool) (int, error) {
	if dryRun {
		log.Println("DRY RUN MODE: Showing what would be posted")
	} else {
		log.Println("posting bookmarks to " + mastodonCredentials.serverDomain)
	}

	var postCount int
	for i := 0; i < len(popular); i++ {
		found, err := db.FindBookmark(popular[i].BookmarkId)
		if err != nil {
			log.Println("error finding bookmark in store: " + popular[i].BookmarkId)
			return postCount, err
		}
		if !found {
			if dryRun {
				log.Printf("DRY RUN: Would post new bookmark:\nID: %s\nTitle: %s\nURL: %s\n",
					popular[i].BookmarkId, popular[i].Title, popular[i].Url)
			} else {
				db.StoreBookmark(*popular[i])
				log.Println("new bookmark stored")
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
	popular, err := ScrapePinboardPopular()
	if err != nil {
		log.Println("error scraping pinboard popular")
		return nil, err
	}

	log.Println("current popular bookmarks: ")
	for i := range popular {
		log.Println(popular[i].BookmarkId)
		log.Println(popular[i].Title)
		log.Println(popular[i].Url)
		log.Println()
	}

	log.Println("found " + fmt.Sprint(len(popular)) + " bookmarks on pinboard popular")
	return popular, err //TODO handle any errors
}
