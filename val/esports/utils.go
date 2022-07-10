package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/val"
)

//go:embed data.json
var parametersFile []byte

type esportsParameters struct {
	Locale string `json:"locale"`
	Title  string `json:"title"`
}

func valEsportsEntryToFeedEntry(entry val.EsportsEntry) internal.FeedEntry {
	return internal.FeedEntry{
		Title:     entry.Title,
		Summary:   entry.Description,
		Authors:   entry.Authors,
		Link:      entry.URL,
		Image:     entry.Image,
		CreatedAt: entry.Date,
		UpdatedAt: entry.Date,
	}
}

func createValEsportsFeed(parameters esportsParameters, entries []val.EsportsEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = valEsportsEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: fmt.Sprintf("https://valorantesports.com/news"),
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

func compareValEsportsEntry(a val.EsportsEntry, b val.EsportsEntry) bool {
	return a.UID == b.UID
}
