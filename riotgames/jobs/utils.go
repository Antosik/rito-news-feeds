package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/riotgames"
)

//go:embed data.json
var parametersFile []byte

type jobsParameters struct {
	Locale string `json:"locale"`
	Title  string `json:"title"`
}

func riotgamesJobsEntryToFeedEntry(entry riotgames.JobsEntry) internal.FeedEntry {
	return internal.FeedEntry{
		Title:      entry.Title,
		Categories: []string{entry.Craft.Name, entry.Products, entry.Office.Name},
		Link:       entry.URL,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func createRiotGamesJobsFeed(parameters jobsParameters, entries []riotgames.JobsEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = riotgamesJobsEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: "https://riotgames.com/" + strings.ToLower(parameters.Locale),
	}

	return internal.Feed{
		Title:    parameters.Title,
		Links:    links,
		Language: parameters.Locale,
		TTL:      uint8(ttl),
		Items:    feedEntries,
	}
}

func getJobsParameters() ([]jobsParameters, error) {
	var data []jobsParameters
	if err := json.Unmarshal(parametersFile, &data); err != nil {
		return nil, fmt.Errorf("can't parse data file: %w", err)
	}

	return data, nil
}

func getRiotGamesJobsEntryKey(item riotgames.JobsEntry) string {
	return item.URL
}
