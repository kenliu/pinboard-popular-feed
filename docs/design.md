# High level design

- Written in go
- Open source
- Web scraper
- DB backend for persistence

## Basic operation
- Cron job for periodically checking the latest posts
- Cron job scrapes the website popular page
- Initially scrapes only the first page
- For each link, check in the DB to see if the ID already exists
- If the ID doesn't exist, post the link to Mastodon, then store it in the DB

# Components

## Web Scraping
Use the `colly` library to scrape the Pinboard popular page and fetch all of the current popular links.

## Datastore
TBD

## Mastodon client
TBD

# CI, deployment, and configuration
TBD

## Cron runner
TBD

# Open Problems
- [ ] rate limiting for posting to mastodon?
- [ ] error handling and monitoring
- [ ] logging
