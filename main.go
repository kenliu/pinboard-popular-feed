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

	// set up log file
	logFile, err := os.OpenFile("pinboard-popular-feed.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("error opening log file")
		os.Exit(1)
	}
	defer logFile.Close()
	log.SetOutput(logFile)

	// set up mastodon credentials
	mastodonCredentials, err := buildMastodonCredentials()
	if err != nil {
		os.Exit(1)
	}

	popular, err := fetchCurrentPinboardPopular()
	if err != nil {
		log.Println("error scraping popular bookmarks")
		os.Exit(1)
	}

	// initialize the bookmark store
	db := data.Init()
	db.InitStore(data.DBConfig{})

	postCount, err := postNewLinks(popular, db, mastodonCredentials, *dryRun)
	if err != nil {
		log.Println("error posting new links")
		os.Exit(1)
	}
	log.Println("posted " + fmt.Sprint(postCount) + " new bookmarks")
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
	log.Println("dryrun: not posting any new bookmarks")
	var postCount int
	for i := 0; i < len(popular); i++ {
		found, err := db.FindBookmark(popular[i].Id)
		if err != nil {
			log.Println("error finding bookmark in store: " + popular[i].Id)
			return postCount, err
		}
		if !found {
			if dryRun {
				log.Println("dry run: new bookmark found, but not posted")
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
	popular, _ := ScrapePinboardPopular()
	log.Println("current popular bookmarks: ")
	for i := 0; i < len(popular); i++ {
		log.Println(popular[i].Id)
		log.Println(popular[i].Title)
		log.Println(popular[i].Url)
		log.Println()
	}

	log.Println("found " + fmt.Sprint(len(popular)) + " bookmarks on pinboard popular")
	return popular, nil //TODO handle any errors
}
