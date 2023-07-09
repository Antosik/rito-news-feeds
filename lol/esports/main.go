package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lol"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	params    []esportsParameters
	paramsErr error

	domain string

	esportsProcessor LoLEsportsProcessor
	mainProcessor    internal.MainProcessor[esportsParameters]
)

const (
	articlesCount = 100
	channelsCount = 5
)

func init() {
	params, paramsErr = getEsportsParameters()
	domain = os.Getenv("DOMAIN_NAME")

	esportsProcessor = LoLEsportsProcessor{}
	mainProcessor = internal.MainProcessor[esportsParameters]{
		Name:        "lolesports",
		Concurrency: channelsCount,

		TypeProcessor: &esportsProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// LoL Esports Processor (implements AbstractProcessor)
type LoLEsportsProcessor struct{}

func (p *LoLEsportsProcessor) GenerateFilePath(param esportsParameters) string {
	return internal.FormatFilePath(filepath.Join("lol", param.Locale, "esports"))
}

func (p *LoLEsportsProcessor) GenerateAbstractFilePath(param esportsParameters) string {
	return fmt.Sprintf("/%s.*", p.GenerateFilePath(param))
}

func (p *LoLEsportsProcessor) ProcessParameters(
	param esportsParameters,
) ([]internal.FeedFile, []error) {
	var (
		client          = lol.EsportsClient{Locale: strings.ToLower(param.Locale)}
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

	existingEntries, err := internal.GetExistingRawEntries[lol.EsportsEntry](existingFile)
	if err != nil {
		errorsCollector.Collect(err)
	}

	diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getLolEsportsEntryKey)
	if isEqual {
		fmt.Printf("%s doesn't require update\n", param.Locale)
		return nil, nil
	}

	fmt.Printf("Found diff: %s...\n", diff)
	fmt.Printf("Updating %s...\n", param.Locale)

	// Create Feed
	feed := createLolEsportsFeed(param, entries)

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

// #region Lambda handler
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
