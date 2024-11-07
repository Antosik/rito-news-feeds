package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/val"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	params    []newsParameters
	errParams error

	domain string

	newsProcessor VALNewsProcessor
	mainProcessor internal.MainProcessor[newsParameters]
)

const (
	articlesCount = 100
	channelsCount = 10
)

func init() {
	params, errParams = getNewsParameters()
	domain = os.Getenv("DOMAIN_NAME")

	newsProcessor = VALNewsProcessor{}
	mainProcessor = internal.MainProcessor[newsParameters]{
		Name:        "valnews",
		Concurrency: channelsCount,

		TypeProcessor: &newsProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// VAL News Processor (implements AbstractProcessor).
type VALNewsProcessor struct{}

func (p *VALNewsProcessor) GenerateFilePath(param newsParameters) string {
	return internal.FormatFilePath(filepath.Join("val", param.Locale, "news"))
}

func (p *VALNewsProcessor) GenerateAbstractFilePath(param newsParameters) string {
	return internal.FormatAbstractFilePath(p.GenerateFilePath(param))
}

func (p *VALNewsProcessor) ProcessParameters(
	param newsParameters,
) ([]internal.FeedFile, []error) {
	var (
		client          = val.NewsClient{Locale: strings.ToLower(param.Locale)}
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
	rawpath := fpath + ".json"

	existingFile, err := mainProcessor.S3Client.DownloadFile(rawpath)
	if err != nil {
		errorsCollector.Collect(err)
	}

	existingEntries, err := internal.GetExistingRawEntries[val.NewsEntry](existingFile)
	if err != nil {
		errorsCollector.Collect(err)
	}

	diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getValNewsEntryKey)
	if isEqual {
		//nolint:forbidigo // need for lambda logs
		fmt.Printf("%s doesn't require update\n", param.Locale)

		return nil, nil
	}

	fmt.Printf("Found diff: %s...\n", diff)      //nolint:forbidigo // need for lambda logs
	fmt.Printf("Updating %s...\n", param.Locale) //nolint:forbidigo // need for lambda logs

	// Create Feed
	feed := createValNewsFeed(param, entries)

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
		return fmt.Errorf("no params found: %w", errParams)
	}

	if domain == "" {
		return internal.ErrDomainNotFound
	}

	return mainProcessor.Process(params)
}

func main() {
	lambda.Start(handler)
}
