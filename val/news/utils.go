package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/val"
)

//go:embed data.json
var parametersFile []byte

type newsParameters struct {
	Locale      string `json:"locale"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func valNewsEntryToFeedEntry(entry val.NewsEntry) internal.FeedEntry {
	return internal.FeedEntry{
		Title:     entry.Title,
		Summary:   entry.Description,
		Link:      entry.URL,
		Image:     entry.Image,
		CreatedAt: entry.Date,
		UpdatedAt: entry.Date,
	}
}

func createValNewsFeed(parameters newsParameters, entries []val.NewsEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = valNewsEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: fmt.Sprintf("https://playvalorant.com/%s/news/", strings.ToLower(parameters.Locale)),
	}

	return internal.Feed{
		Title:       parameters.Title,
		Description: parameters.Description,
		Links:       links,
		Language:    parameters.Locale,
		TTL:         uint8(ttl),
		Items:       feedEntries,
	}
}

func getNewsParameters() ([]newsParameters, error) {
	var data []newsParameters
	if err := json.Unmarshal(parametersFile, &data); err != nil {
		return nil, fmt.Errorf("can't parse data file: %w", err)
	}

	return data, nil
}

func getValNewsEntryKey(item val.NewsEntry) string {
	return item.URL
}
