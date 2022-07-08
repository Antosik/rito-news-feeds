package internal

import (
	"time"
)

// JSON Feed Docs - https://jsonfeed.org/version/1

type JSONFeed struct {
	Version     string          `json:"version"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	HomePageUrl string          `json:"home_page_url"`
	FeedUrl     string          `json:"feed_url"`
	Language    string          `json:"language"`
	Author      *JSONFeedAuthor `json:"author,omitempty"`
	Items       []JSONFeedEntry `json:"items"`
}

func (jsonFeed JSONFeed) JSON() ([]byte, error) {
	return MarshalJSON(jsonFeed)
}

type JSONFeedAuthor struct {
	Name string `json:"name"`
	URL  string `json:"url,omitempty"`
}

type JSONFeedEntry struct {
	ID          string           `json:"id"`
	Title       string           `json:"title"`
	ContentText string           `json:"content_text,omitempty"`
	Image       string           `json:"image,omitempty"`
	URL         string           `json:"url"`
	Author      []JSONFeedAuthor `json:"authors,omitempty"`
	Tags        []string         `json:"tags,omitempty"`
	Published   string           `json:"date_published"`
}

func ConvertFeedEntryToJSONFeedEntry(entry *FeedEntry) JSONFeedEntry {
	var (
		authors = make([]JSONFeedAuthor, len(entry.Authors))
		tags    = make([]string, len(entry.Categories))
	)

	for i, author := range entry.Authors {
		authors[i] = JSONFeedAuthor{Name: author}
	}

	for i, category := range entry.Categories {
		tags[i] = category
	}

	return JSONFeedEntry{
		ID:          entry.Link,
		Title:       entry.Title,
		ContentText: entry.Summary,
		Image:       entry.Image,
		URL:         entry.Link,
		Author:      authors,
		Tags:        tags,
		Published:   entry.CreatedAt.Format(time.RFC3339),
	}
}

func ConvertFeedToJSONFeed(feed *Feed) JSONFeed {
	entries := make([]JSONFeedEntry, len(feed.Items))
	for i, item := range feed.Items {
		entries[i] = item.JSONFeed()
	}

	return JSONFeed{
		Version:     "https://jsonfeed.org/version/1.1",
		Title:       feed.Title,
		Description: feed.Description,
		HomePageUrl: feed.Links.Alternate,
		FeedUrl:     feed.Links.Self,
		Language:    feed.Language,
		Author: &JSONFeedAuthor{
			Name: "Antosik",
			URL:  "https://iamantosik.me",
		},
		Items: entries,
	}
}
