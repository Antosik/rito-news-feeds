package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lol"
)

//go:embed data.json
var parametersFile []byte

type esportsParameters struct {
	Locale string `json:"locale"`
	Title  string `json:"title"`
}

func lolEsportsEntryToFeedEntry(entry lol.EsportsEntry) internal.FeedEntry {
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

func createLolEsportsFeed(parameters esportsParameters, entries []lol.EsportsEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = lolEsportsEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: "https://lolesports.com/news",
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

func compareLolEsportsEntry(a lol.EsportsEntry, b lol.EsportsEntry) bool {
	return a.UID == b.UID
}
