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

func GenerateFeedFiles(feed Feed, domain string, name string) ([]FeedFile, []error) {
	var (
		generatedFiles   = make([]FeedFile, 0, 4)
		generationErrors = make([]error, 0, 4)
		formats          = []string{"atom", "rss", "jsonfeed"}
	)

	for _, format := range formats {
		mime := feedMimeType[format]
		filename := fmt.Sprintf("%s.%s", name, format)

		if domain != "" {
			feed.Links.Self = fmt.Sprintf("https://%s/%s", domain, filename)
		}

		data, err := feedBufferGenerators[format](feed)
		if err != nil {
			generationErrors = append(generationErrors, err)
		} else {
			generatedFiles = append(generatedFiles, FeedFile{
				Name:     filename,
				MimeType: mime,
				Buffer:   data,
			})
		}
	}

	return generatedFiles, generationErrors
}

func GenerateRawFile(entries interface{}, name string) (FeedFile, error) {
	rawjson, err := MarshalJSON(entries)
	if err != nil {
		return FeedFile{}, err
	}

	return FeedFile{
		Name:     name,
		MimeType: "application/json",
		Buffer:   rawjson,
	}, nil
}
