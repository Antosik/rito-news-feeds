package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/riotgames"
)

//go:embed data.json
var parametersFile []byte

type newsParameters struct {
	Locale      string `json:"locale"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func riotGamesNewsEntryToFeedEntry(entry riotgames.NewsEntry) internal.FeedEntry {
	return internal.FeedEntry{
		Title:      entry.Title,
		Summary:    entry.Description,
		Categories: []string{entry.Category},
		Link:       entry.URL,
		Image:      entry.Image,
		CreatedAt:  entry.Date,
		UpdatedAt:  entry.Date,
	}
}

func createRiotGamesNewsFeed(parameters newsParameters, entries []riotgames.NewsEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = riotGamesNewsEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: fmt.Sprintf("https://www.riotgames.com/%s", strings.ToLower(parameters.Locale)),
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

func compareRiotGamesNewsEntry(a riotgames.NewsEntry, b riotgames.NewsEntry) bool {
	return a.URL == b.URL
}
