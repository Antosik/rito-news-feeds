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
	params    []esportsParameters
	errParams error

	domain string

	esportsProcessor VALEsportsProcessor
	mainProcessor    internal.MainProcessor[esportsParameters]
)

const (
	articlesCount = 100
	channelsCount = 10
)

func init() {
	params, errParams = getEsportsParameters()
	domain = os.Getenv("DOMAIN_NAME")

	esportsProcessor = VALEsportsProcessor{}
	mainProcessor = internal.MainProcessor[esportsParameters]{
		Name:        "valesports",
		Concurrency: channelsCount,

		TypeProcessor: &esportsProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// VAL Esports Processor (implements AbstractProcessor).
type VALEsportsProcessor struct{}

func (p *VALEsportsProcessor) GenerateFilePath(param esportsParameters) string {
	return internal.FormatFilePath(filepath.Join("val", param.Locale, "esports"))
}

func (p *VALEsportsProcessor) GenerateAbstractFilePath(param esportsParameters) string {
	return internal.FormatAbstractFilePath(p.GenerateFilePath(param))
}

func (p *VALEsportsProcessor) ProcessParameters(
	param esportsParameters,
) ([]internal.FeedFile, []error) {
	var (
		client          = val.EsportsClient{Locale: strings.ToLower(param.Locale)}
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

	existingEntries, err := internal.GetExistingRawEntries[val.EsportsEntry](existingFile)
	if err != nil {
		errorsCollector.Collect(err)
	}

	diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getValEsportsEntryKey)
	if isEqual {
		//nolint:forbidigo // need for lambda logs
		fmt.Printf("%s doesn't require update\n", param.Locale)

		return nil, nil
	}

	fmt.Printf("Found diff: %s...\n", diff)      //nolint:forbidigo // need for lambda logs
	fmt.Printf("Updating %s...\n", param.Locale) //nolint:forbidigo // need for lambda logs

	// Create Feed
	feed := createValEsportsFeed(param, entries)

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
