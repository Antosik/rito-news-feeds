package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/riotgames"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	params    []newsParameters
	errParams error

	domain string

	newsProcessor RiotGamesNewsProcessor
	mainProcessor internal.MainProcessor[newsParameters]
)

const (
	articlesCount = 100
	channelsCount = 10
)

func init() {
	params, errParams = getNewsParameters()
	domain = os.Getenv("DOMAIN_NAME")

	newsProcessor = RiotGamesNewsProcessor{}
	mainProcessor = internal.MainProcessor[newsParameters]{
		Name:        "riotgamesnews",
		Concurrency: channelsCount,

		TypeProcessor: &newsProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// RiotGames News Processor (implements AbstractProcessor).
type RiotGamesNewsProcessor struct{}

func (p *RiotGamesNewsProcessor) GenerateFilePath(param newsParameters) string {
	return internal.FormatFilePath(filepath.Join("riotgames", param.Locale, "news"))
}

func (p *RiotGamesNewsProcessor) GenerateAbstractFilePath(param newsParameters) string {
	return internal.FormatAbstractFilePath(p.GenerateFilePath(param))
}

func (p *RiotGamesNewsProcessor) ProcessParameters(
	param newsParameters,
) ([]internal.FeedFile, []error) {
	var (
		client          = riotgames.NewsClient{Locale: strings.ToLower(param.Locale)}
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

	existingEntries, err := internal.GetExistingRawEntries[riotgames.NewsEntry](existingFile)
	if err != nil {
		errorsCollector.Collect(err)
	}

	diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getRiotGamesNewsEntryKey)
	if isEqual {
		//nolint:forbidigo // need for lambda logs
		fmt.Printf("%s doesn't require update\n", param.Locale)

		return nil, nil
	}

	fmt.Printf("Found diff: %s...\n", diff)      //nolint:forbidigo // need for lambda logs
	fmt.Printf("Updating %s...\n", param.Locale) //nolint:forbidigo // need for lambda logs

	// Create Feed
	feed := createRiotGamesNewsFeed(param, entries)

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
