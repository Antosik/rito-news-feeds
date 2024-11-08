package internal

import (
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

// RSS Docs - http://cyber.harvard.edu/rss/rss.html

const NSDublinCore = "http://purl.org/dc/elements/1.1/"

type RSS struct {
	XMLName      string     `xml:"rss"`
	Version      string     `xml:"version,attr"`
	NSAtom       string     `xml:"xmlns:atom,attr"`
	NSDublinCore string     `xml:"xmlns:dc,attr"`
	Channel      RSSChannel `xml:"channel"`
}

func (rss RSS) XML() ([]byte, error) {
	data, err := xml.Marshal(rss)
	if err != nil {
		return nil, fmt.Errorf("can't marshal rss feed: %w", err)
	}

	return data, nil
}

type RSSChannel struct {
	Title         string     `xml:"title"`
	Link          string     `xml:"link"`
	Description   string     `xml:"description"`
	Language      string     `xml:"language"`
	LastBuildDate string     `xml:"lastBuildDate"`
	Generator     string     `xml:"generator"`
	Docs          string     `xml:"docs"`
	TTL           uint8      `xml:"ttl"`
	AtomLink      AtomLink   `xml:"atom:link"`
	Item          []RSSEntry `xml:"item"`
}

type RSSEnclosure struct {
	URL    string `xml:"url,attr"`
	Length int    `xml:"length,attr"`
	Type   string `xml:"type,attr"`
}

type RSSEntry struct {
	Title       string        `xml:"title"`
	Description string        `xml:"description,omitempty"`
	Published   string        `xml:"pubDate"`
	Link        string        `xml:"link"`
	GUID        string        `xml:"guid"`
	Enclosure   *RSSEnclosure `xml:"enclosure,omitempty"`
	Author      string        `xml:"dc:creator,omitempty"`
	Category    []string      `xml:"category,omitempty"`
}

func ConvertFeedEntryToRSSEntry(entry *FeedEntry) RSSEntry {
	var image *RSSEnclosure

	if entry.Image != "" {
		image = &RSSEnclosure{URL: entry.Image, Length: 0, Type: "image/*"}
	}

	return RSSEntry{
		Title:       entry.Title,
		Description: entry.Summary,
		Published:   entry.CreatedAt.Format(time.RFC1123Z),
		Link:        entry.Link,
		GUID:        entry.Link,
		Enclosure:   image,
		Author:      strings.Join(entry.Authors, ", "),
		Category:    entry.Categories,
	}
}

func ConvertFeedToRSS(feed *Feed) RSS {
	entries := make([]RSSEntry, len(feed.Items))
	for i, item := range feed.Items {
		entries[i] = item.RSS()
	}

	return RSS{
		Version:      "2.0",
		NSAtom:       NSAtom,
		NSDublinCore: NSDublinCore,
		Channel: RSSChannel{
			Title:         feed.Title,
			Link:          feed.Links.Alternate,
			Description:   feed.Description,
			Language:      feed.Language,
			LastBuildDate: time.Now().UTC().Format(time.RFC1123Z),
			Generator:     "rito-news-feeds",
			Docs:          "https://github.com/Antosik/rito-news-feeds",
			TTL:           feed.TTL,
			Item:          entries,
			AtomLink:      AtomLink{feed.Links.Self, "self", "application/rss+xml"},
		},
	}
}
