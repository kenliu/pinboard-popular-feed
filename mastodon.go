package main

import (
	"log"
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

	log.Println("posting to mastodon: " + tootText)
	client := http.Client{}
	endpoint := "https://" + credentials.serverDomain + "/api/v1/statuses"

	form := url.Values{}
	form.Add("status", tootText)

	req, err := http.NewRequest("POST", endpoint, strings.NewReader(form.Encode()))
	//generate an error intentionally
	//req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		log.Println("error creating request")
		return err
	}

	req.Header.Add("Authorization", "Bearer "+credentials.accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error posting to mastodon")
		return err
	}
	if resp.Status != "200 OK" {
		log.Println("error posting to mastodon: " + resp.Status)
		log.Println(resp)
		return err
	}

	// TODO log this as a debug message
	log.Println("posted to mastodon")
	return nil
}

func buildToot(b data.Bookmark) string {
	toot := b.Title + "\n" + b.Url
	return toot
}
