package internal

import (
	"time"
)

type FeedEntry struct {
	Title      string
	Summary    string
	Image      string
	Link       string
	Authors    []string
	Categories []string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (entry *FeedEntry) Atom() AtomEntry {
	return ConvertFeedEntryToAtomEntry(entry)
}

func (entry *FeedEntry) RSS() RSSEntry {
	return ConvertFeedEntryToRSSEntry(entry)
}

func (entry *FeedEntry) JSONFeed() JSONFeedEntry {
	return ConvertFeedEntryToJSONFeedEntry(entry)
}

type FeedLinks struct {
	Self      string
	Alternate string
}

type Feed struct {
	Title       string
	Description string
	Links       FeedLinks
	Language    string
	TTL         uint8
	Items       []FeedEntry
}

func (feed *Feed) Atom() Atom {
	return ConvertFeedToAtom(feed)
}

func (feed *Feed) RSS() RSS {
	return ConvertFeedToRSS(feed)
}

func (feed *Feed) JSONFeed() JSONFeed {
	return ConvertFeedToJSONFeed(feed)
}
