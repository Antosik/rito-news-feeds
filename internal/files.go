package internal

import (
	"fmt"
)

type FeedFile struct {
	Name     string
	MimeType string
	Buffer   []byte
}

var feedMimeType = map[string]string{
	"atom":     "application/atom+xml",
	"rss":      "application/rss+xml",
	"jsonfeed": "application/feed+json",
}

var feedBufferGenerators = map[string]func(feed Feed) ([]byte, error){
	"atom":     func(feed Feed) ([]byte, error) { return feed.Atom().XML() },
	"rss":      func(feed Feed) ([]byte, error) { return feed.RSS().XML() },
	"jsonfeed": func(feed Feed) ([]byte, error) { return feed.JSONFeed().JSON() },
}

func GenerateFeedFiles(feed Feed, name string) ([]FeedFile, []error) {
	var (
		generatedFiles   = make([]FeedFile, 0, 4)
		generationErrors = make([]error, 0, 4)
		formats          = []string{"atom", "rss", "jsonfeed"}
	)

	for _, format := range formats {
		mime := feedMimeType[format]

		data, err := feedBufferGenerators[format](feed)
		if err != nil {
			generationErrors = append(generationErrors, err)
		} else {
			filename := fmt.Sprintf("%s.%s", name, format)
			generatedFiles = append(generatedFiles, FeedFile{
				Name:     filename,
				MimeType: mime,
				Buffer:   data,
			})
		}
	}

	return generatedFiles, generationErrors
}
