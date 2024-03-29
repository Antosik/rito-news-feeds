package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lol"
)

//go:embed data.json
var parametersFile []byte

//go:embed locales.json
var localeFile []byte

type statusParameters struct {
	ID      string   `json:"id"`
	Region  string   `json:"region"`
	Locales []string `json:"locales"`
}

type statusLocale struct {
	Locale string `json:"locale"`
	Title  string `json:"title"`
}

func lolStatusEntryToFeedEntry(entry lol.StatusEntry) internal.FeedEntry {
	return internal.FeedEntry{
		Title:     entry.Title,
		Summary:   entry.Description,
		Authors:   []string{entry.Author},
		Link:      entry.URL,
		CreatedAt: entry.Date,
		UpdatedAt: entry.Date,
	}
}

func createLolStatusFeed(regionID string, locale *statusLocale, entries []lol.StatusEntry) internal.Feed {
	feedEntries := make([]internal.FeedEntry, len(entries))
	for i, entry := range entries {
		feedEntries[i] = lolStatusEntryToFeedEntry(entry)
	}

	ttl, err := strconv.ParseUint(os.Getenv("TTL"), 10, 8)
	if err != nil {
		ttl = 15
	}

	links := internal.FeedLinks{
		Alternate: fmt.Sprintf("https://status.riotgames.com/lol?region=%s&locale=%s", regionID, locale.Locale),
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

func getLolStatusEntryKey(item lol.StatusEntry) string {
	return item.URL
}
