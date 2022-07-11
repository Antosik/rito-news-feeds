package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/wr"
)

//go:embed data.json
var parametersFile []byte

type esportsParameters struct {
	Locale string `json:"locale"`
	Title  string `json:"title"`
}

func wrEsportsEntryToFeedEntry(entry wr.EsportsEntry) internal.FeedEntry {
	return internal.FeedEntry{
		Title:      entry.Title,
		Summary:    entry.Description,
		Authors:    entry.Authors,
		Categories: entry.Categories,
		Link:       entry.URL,
		Image:      entry.Image,
		CreatedAt:  entry.Date,
		UpdatedAt:  entry.Date,
	}
}

func createWrEsportsFeed(parameters esportsParameters, entries []wr.EsportsEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = wrEsportsEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: fmt.Sprintf("https://wildriftesports.com/%s/news", parameters.Locale),
	}

	return internal.Feed{
		Title:    parameters.Title,
		Links:    links,
		Language: parameters.Locale,
		TTL:      uint8(ttl),
		Items:    feedEntries,
	}
}

func getEsportsParameters() ([]esportsParameters, error) {
	var data []esportsParameters
	if err := json.Unmarshal(parametersFile, &data); err != nil {
		return nil, fmt.Errorf("can't parse data file: %w", err)
	}

	return data, nil
}

func compareWrEsportsEntry(a wr.EsportsEntry, b wr.EsportsEntry) bool {
	return a.UID == b.UID
}
