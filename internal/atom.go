package internal

import (
	"encoding/xml"
	"strings"
	"time"
)

const NSAtom = "http://www.w3.org/2005/Atom"

// Atom Docs - https://validator.w3.org/feed/docs/atom.html

type Atom struct {
	XMLName   string        `xml:"feed"`
	Xmlns     string        `xml:"xmlns,attr"`
	Lang      string        `xml:"xml:lang,attr"`
	ID        string        `xml:"id"`
	Title     string        `xml:"title"`
	Subtitle  string        `xml:"subtitle,omitempty"`
	Author    AtomAuthor    `xml:"author"`
	Generator AtomGenerator `xml:"generator"`
	Link      []AtomLink    `xml:"link"`
	Updated   string        `xml:"updated"`
	Entries   []AtomEntry   `xml:"entry"`
}

func (atom Atom) XML() ([]byte, error) {
	return xml.Marshal(atom)
}

type AtomGenerator struct {
	Name    string `xml:",chardata"`
	URI     string `xml:"uri,attr"`
	Version string `xml:"version,attr,omitempty"`
}

type AtomEntry struct {
	ID        string         `xml:"id"`
	Title     string         `xml:"title"`
	Summary   string         `xml:"summary,omitempty"`
	Link      []AtomLink     `xml:"link"`
	Author    []AtomAuthor   `xml:"author,omitempty"`
	Category  []AtomCategory `xml:"category,omitempty"`
	Published string         `xml:"published"`
	Updated   string         `xml:"updated"`
}

type AtomAuthor struct {
	Name string `xml:"name"`
	URI  string `xml:"uri,omitempty"`
}

type AtomContent struct {
	Type  string `xml:"type,attr"`
	Value string `xml:",cdata"`
}

type AtomCategory struct {
	Term  string `xml:"term,attr"`
	Label string `xml:"label,attr"`
}

type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

func ConvertFeedEntryToAtomEntry(entry *FeedEntry) AtomEntry {
	var (
		links      = make([]AtomLink, 0, 2)
		authors    = make([]AtomAuthor, len(entry.Authors))
		categories = make([]AtomCategory, len(entry.Categories))
	)

	links = append(links, AtomLink{Href: entry.Link, Rel: "alternate"})
	if entry.Image != "" {
		links = append(links, AtomLink{Href: entry.Image, Rel: "enclosure"})
	}

	for i, author := range entry.Authors {
		authors[i] = AtomAuthor{Name: author}
	}

	for i, category := range entry.Categories {
		categories[i] = AtomCategory{Label: category, Term: strings.ToLower(category)}
	}

	return AtomEntry{
		ID:        entry.Link,
		Title:     entry.Title,
		Summary:   entry.Summary,
		Link:      links,
		Author:    authors,
		Category:  categories,
		Published: entry.CreatedAt.Format(time.RFC3339),
		Updated:   entry.UpdatedAt.Format(time.RFC3339),
	}
}

func ConvertFeedToAtom(feed *Feed) Atom {
	entries := make([]AtomEntry, len(feed.Items))
	for i, item := range feed.Items {
		entries[i] = item.Atom()
	}

	links := make([]AtomLink, 0, 2)

	links = append(links, AtomLink{
		Href: feed.Links.Alternate,
		Rel:  "alternate",
	})

	if feed.Links.Self != "" {
		links = append(links, AtomLink{
			Href: feed.Links.Self,
			Rel:  "self",
			Type: "application/atom+xml",
		})
	}

	return Atom{
		Xmlns:    NSAtom,
		Lang:     feed.Language,
		ID:       feed.Links.Self,
		Title:    feed.Title,
		Subtitle: feed.Description,
		Author: AtomAuthor{
			Name: "Antosik",
			URI:  "https://iamantosik.me",
		},
		Generator: AtomGenerator{
			Name: "rito-news-feeds",
			URI:  "https://github.com/Antosik/rito-news-feeds",
		},
		Link:    links,
		Updated: time.Now().UTC().Format(time.RFC3339),
		Entries: entries,
	}
}
