package main

import (
	"log/slog"
	"net/http"
	"net/url"
	"pinboard-popular-feed/data"
	"strings"
)

type MastodonCredentials struct {
	serverDomain string
	accessToken  string
}

func TootBookmark(b data.Bookmark, credentials MastodonCredentials) error {
	tootText := buildToot(b)

	slog.Info("posting to mastodon", "toot", tootText)
	client := http.Client{}
	endpoint := "https://" + credentials.serverDomain + "/api/v1/statuses"

	form := url.Values{}
	form.Add("status", tootText)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	//generate an error intentionally
	//req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		slog.Error("error creating request", "error", err)
		return err
	}

	req.Header.Add("Authorization", "Bearer "+credentials.accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		slog.Error("error posting to mastodon", "error", err)
		return err
	}
	if resp.Status != "200 OK" {
		slog.Error("error posting to mastodon", "status", resp.Status)
		return err
	}

	slog.Debug("successfully posted to mastodon")
	return nil
}

func buildToot(b data.Bookmark) string {
	toot := b.Title + "\n" + b.Url
	return toot
}
