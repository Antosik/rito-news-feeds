package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lor"
)

//go:embed data.json
var parametersFile []byte

//go:embed locales.json
var localeFile []byte

type statusParameters struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type statusLocale struct {
	Locale string `json:"locale"`
	Title  string `json:"title"`
}

func lorStatusEntryToFeedEntry(entry lor.StatusEntry) internal.FeedEntry {
	return internal.FeedEntry{
		Title:     entry.Title,
		Summary:   entry.Description,
		Authors:   []string{entry.Author},
		Link:      entry.URL,
		CreatedAt: entry.Date,
		UpdatedAt: entry.Date,
	}
}

func createLorStatusFeed(regionID string, locale statusLocale, entries []lor.StatusEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = lorStatusEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: fmt.Sprintf("https://status.riotgames.com/lor?region=%s&locale=%s", regionID, locale.Locale),
	}

	return internal.Feed{
		Title:    locale.Title,
		Links:    links,
		Language: strings.ReplaceAll(locale.Locale, "_", "-"),
		TTL:      uint8(ttl),
		Items:    feedEntries,
	}
}

func getStatusParameters() ([]statusParameters, error) {
	var data []statusParameters
	if err := json.Unmarshal(parametersFile, &data); err != nil {
		return nil, fmt.Errorf("can't parse data file: %w", err)
	}

	return data, nil
}

func getStatusLocales() (map[string]statusLocale, error) {
	var data []statusLocale
	if err := json.Unmarshal(localeFile, &data); err != nil {
		return nil, fmt.Errorf("can't parse data file: %w", err)
	}

	localeMap := make(map[string]statusLocale, len(data))
	for _, locale := range data {
		localeMap[locale.Locale] = locale
	}

	return localeMap, nil
}

func getLorStatusEntryKey(item lor.StatusEntry) string {
	return item.URL
}
