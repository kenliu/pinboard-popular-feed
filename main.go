package main

import (
	"errors"
	"flag"
	"log/slog"
	"os"
	"pinboard-popular-feed/data"
)

// Add these functions before main()
func isCloudRun() bool {
	return os.Getenv("CLOUD_RUN_JOB") != ""
}

func setupLogger() {
	var handler slog.Handler
	if isCloudRun() {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}
	logger := slog.New(handler)
	slog.SetDefault(logger)
}

func main() {
	// handle dry-run command line flag
	dryRun := flag.Bool("dryrun", false, "scan for new posts but don't post to mastodon")
	flag.Parse()

	// set up logging
	setupLogger()

	slog.Info("starting pinboard-popular-feed processing")
	if *dryRun {
		slog.Info("DRY RUN MODE: showing what would be posted")
	}

	// set up mastodon credentials
	mastodonCredentials, err := buildMastodonCredentials()
	if err != nil {
		slog.Error("error setting up mastodon credentials", "error", err)
		os.Exit(1)
	}

	popular, err := fetchCurrentPinboardPopular()
	if err != nil {
		slog.Error("error scraping popular bookmarks", "error", err)
		os.Exit(1)
	}

	// initialize the bookmark store
	var db = data.BookmarkStore{}
	db.InitStore(data.CreateDBConfigFromEnv())

	_, err = postNewLinks(popular, db, mastodonCredentials, *dryRun)
	if err != nil {
		slog.Error("error posting new bookmarks", "error", err)
		os.Exit(1)
	}

	slog.Info("finished pinboard-popular-feed processing")
}

func buildMastodonCredentials() (MastodonCredentials, error) {
	if os.Getenv("MASTODON_ACCESS_TOKEN") == "" {
		slog.Error("MASTODON_ACCESS_TOKEN not set")
		return MastodonCredentials{}, errors.New("MASTODON_ACCESS_TOKEN not set")
	}

	if os.Getenv("MASTODON_SERVER_DOMAIN") == "" {
		slog.Error("MASTODON_SERVER_DOMAIN not set")
		return MastodonCredentials{}, errors.New("MASTODON_SERVER_DOMAIN not set")
	}

	return MastodonCredentials{
		accessToken:  os.Getenv("MASTODON_ACCESS_TOKEN"),
		serverDomain: os.Getenv("MASTODON_SERVER_DOMAIN"),
	}, nil
}

func postNewLinks(popular []*data.Bookmark, db data.BookmarkStore, mastodonCredentials MastodonCredentials, dryRun bool) (int, error) {
	if dryRun {
		slog.Info("DRY RUN MODE: Showing what would be posted")
	} else {
		slog.Info("start posting bookmarks to mastodon", "server", mastodonCredentials.serverDomain)
	}

	var postCount int
	for i := range popular {
		found, err := db.FindBookmark(popular[i].BookmarkId)
		if err != nil {
			slog.Error("error finding bookmark in store", "bookmark_id", popular[i].BookmarkId, "error", err)
			return postCount, err
		}
		if !found {
			if dryRun {
				slog.Info("DRY RUN: Would post new bookmark",
					"id", popular[i].BookmarkId,
					"title", popular[i].Title,
					"url", popular[i].Url)
			} else {
				db.StoreBookmark(*popular[i])
				slog.Info("new bookmark stored", "id", popular[i].BookmarkId)
				postCount++
				TootBookmark(*popular[i], mastodonCredentials)
			}
		}
		// TODO implement some kind of rate limiting?
		// https://docs.joinmastodon.org/api/rate-limits/
	}

	slog.Info("posted all new bookmarks", "count", postCount)
	return postCount, nil
}

func fetchCurrentPinboardPopular() ([]*data.Bookmark, error) {
	popular, err := ScrapePinboardPopular()
	if err != nil {
		slog.Error("error scraping pinboard popular", "error", err)
		return nil, err
	}

	slog.Debug("current popular bookmarks", "count", len(popular))
	for i := range popular {
		slog.Debug("bookmark details",
			"id", popular[i].BookmarkId,
			"title", popular[i].Title,
			"url", popular[i].Url)
	}

	slog.Info("found bookmarks on pinboard popular", "count", len(popular))
	return popular, err
}
