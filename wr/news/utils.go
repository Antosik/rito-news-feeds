package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/wr"
)

//go:embed data.json
var parametersFile []byte

type newsParameters struct {
	Locale      string `json:"locale"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func wrNewsEntryToFeedEntry(entry wr.NewsEntry) internal.FeedEntry {
	categories := make([]string, 0, len(entry.Categories)+len(entry.Tags))
	categories = append(categories, entry.Categories...)
	categories = append(categories, entry.Tags...)

	return internal.FeedEntry{
		Title:      entry.Title,
		Summary:    entry.Description,
		Categories: categories,
		Link:       entry.URL,
		Image:      entry.Image,
		CreatedAt:  entry.Date,
		UpdatedAt:  entry.Date,
	}
}

func createWrNewsFeed(parameters newsParameters, entries []wr.NewsEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = wrNewsEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: fmt.Sprintf("https://wildrift.leagueoflegends.com/%s/news/", strings.ToLower(parameters.Locale)),
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

func compareWrNewsEntry(a wr.NewsEntry, b wr.NewsEntry) bool {
	return a.UID == b.UID
}