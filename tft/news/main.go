package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/tft"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	params    []newsParameters
	paramsErr error

	domain string

	newsProcessor TFTNewsProcessor
	mainProcessor internal.MainProcessor[newsParameters]
)

const (
	articlesCount = 100
	channelsCount = 10
)

func init() {
	params, paramsErr = getNewsParameters()
	domain = os.Getenv("DOMAIN_NAME")

	newsProcessor = TFTNewsProcessor{}
	mainProcessor = internal.MainProcessor[newsParameters]{
		Name:        "tftnews",
		Concurrency: channelsCount,

		TypeProcessor: &newsProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// TFT News Processor (implements AbstractProcessor)
type TFTNewsProcessor struct{}

func (p *TFTNewsProcessor) GenerateFilePath(param newsParameters) string {
	return internal.FormatFilePath(filepath.Join("tft", param.Locale, "news"))
}

func (p *TFTNewsProcessor) GenerateAbstractFilePath(param newsParameters) string {
	return fmt.Sprintf("/%s.*", p.GenerateFilePath(param))
}

func (p *TFTNewsProcessor) ProcessParameters(
	param newsParameters,
) ([]internal.FeedFile, []error) {
	var (
		client          = tft.NewsClient{Locale: strings.ToLower(param.Locale)}
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
	rawpath := fmt.Sprintf("%s.json", fpath)

	existingFile, err := mainProcessor.S3Client.DownloadFile(rawpath)
	if err != nil {
		errorsCollector.Collect(err)
	}

	existingEntries, err := internal.GetExistingRawEntries[tft.NewsEntry](existingFile)
	if err != nil {
		errorsCollector.Collect(err)
	}

	diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getTftNewsEntryKey)
	if isEqual {
		fmt.Printf("%s doesn't require update\n", param.Locale)
		return nil, nil
	}

	fmt.Printf("Found diff: %s...\n", diff)
	fmt.Printf("Updating %s...\n", param.Locale)

	// Create Feed
	feed := createTftNewsFeed(param, entries)

	// Generate Atom, JSONFeed, RSS file
	files, errors := internal.GenerateFeedFiles(feed, domain, fpath)
	if len(errors) > 0 {
		errorsCollector.CollectMany(errors)
	}

	// Generate RAW file
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
