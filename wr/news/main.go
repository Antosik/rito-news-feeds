package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/wr"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	params    []newsParameters
	paramsErr error

	domain string

	newsProcessor WRNewsProcessor
	mainProcessor internal.MainProcessor[newsParameters]
)

const (
	articlesCount = 100
	channelsCount = 8
)

func init() {
	params, paramsErr = getNewsParameters()
	domain = os.Getenv("DOMAIN_NAME")

	newsProcessor = WRNewsProcessor{}
	mainProcessor = internal.MainProcessor[newsParameters]{
		Name:        "wrnews",
		Concurrency: channelsCount,

		TypeProcessor: &newsProcessor,

		CFInvalidator: internal.NewCloudFrontInvalidator(),
		S3Uploader:    internal.NewS3Uploader(),
	}
}

// WR News Processor (implements AbstractProcessor)
type WRNewsProcessor struct{}

func (p *WRNewsProcessor) GenerateFilePath(param newsParameters) string {
	return internal.FormatFilePath(filepath.Join("wr", param.Locale, "news"))
}

func (p *WRNewsProcessor) GenerateInvalidationFilePath(param newsParameters) string {
	return fmt.Sprintf("/%s.*", p.GenerateFilePath(param))
}

func (p *WRNewsProcessor) GenerateAsteriskInvalidationPath() string {
	return filepath.Join("/", "wr", "*", "news.*")
}

func (p *WRNewsProcessor) ProcessParameters(
	param newsParameters,
) ([]internal.FeedFile, []error) {
	var (
		client          = wr.NewsClient{Locale: strings.ToLower(param.Locale)}
		fpath           = p.GenerateFilePath(param)
		errorsCollector = internal.NewErrorCollector()
	)

	// Get new items
	entries, err := client.GetItems(articlesCount)
	if err != nil {
		errorsCollector.Collect(fmt.Errorf("can't get items for %s: %w", param.Locale, err))
		return nil, *errorsCollector
	}

	// Check diff with existing data
	existingEntries, err := internal.GetExistingRawEntries[wr.NewsEntry](domain, fpath)
	if err != nil {
		errorsCollector.Collect(err)
	} else if internal.IsEqual(existingEntries, entries, compareWrNewsEntry) {
		fmt.Printf("%s doesn't require update\n", param.Locale)
		return nil, nil
	}

	fmt.Printf("Updating %s...\n", param.Locale)

	// Create Feed
	feed := createWrNewsFeed(param, entries)

	// Generate Atom, JSONFeed, RSS file
	files, errors := internal.GenerateFeedFiles(feed, domain, fpath)
	if len(errors) > 0 {
		errorsCollector.CollectMany(errors)
	}

	// Generate RAW file
	rawpath := fmt.Sprintf("%s.json", fpath)

	rawfile, err := internal.GenerateRawFile(entries, rawpath)
	if err != nil {
		errorsCollector.Collect(err)
	} else {
		files = append(files, rawfile)
	}

	return files, *errorsCollector
}

func handler() error {
	if len(params) == 0 {
		return fmt.Errorf("no params found: %w", paramsErr)
	}

	if domain == "" {
		return internal.ErrDomainNotFound
	}

	return mainProcessor.Process(params)
}

func main() {
	lambda.Start(handler)
}
