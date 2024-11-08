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
	params    []jobsParameters
	errParams error

	domain string

	newsProcessor RiotGamesJobsProcessor
	mainProcessor internal.MainProcessor[jobsParameters]
)

const (
	channelsCount = 10
)

func init() {
	params, errParams = getJobsParameters()
	domain = os.Getenv("DOMAIN_NAME")

	newsProcessor = RiotGamesJobsProcessor{}
	mainProcessor = internal.MainProcessor[jobsParameters]{
		Name:        "riotgamesjobs",
		Concurrency: channelsCount,

		TypeProcessor: &newsProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// RiotGames Jobs Processor (implements AbstractProcessor).
type RiotGamesJobsProcessor struct{}

func (p *RiotGamesJobsProcessor) GenerateFilePath(param jobsParameters) string {
	return internal.FormatFilePath(filepath.Join("riotgames", param.Locale, "jobs"))
}

func (p *RiotGamesJobsProcessor) GenerateAbstractFilePath(param jobsParameters) string {
	return internal.FormatAbstractFilePath(p.GenerateFilePath(param))
}

func (p *RiotGamesJobsProcessor) ProcessParameters(
	param jobsParameters,
) ([]internal.FeedFile, []error) {
	var (
		client          = riotgames.JobsClient{Locale: strings.ToLower(param.Locale)}
		fpath           = p.GenerateFilePath(param)
		errorsCollector = internal.NewErrorCollector()
	)

	// Get new items
	entries, err := client.GetItems()
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

	existingEntries, err := internal.GetExistingRawEntries[riotgames.JobsEntry](existingFile)
	if err != nil {
		errorsCollector.Collect(err)
	}

	diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getRiotGamesJobsEntryKey)
	if isEqual {
		//nolint:forbidigo // need for lambda logs
		fmt.Printf("%s doesn't require update\n", param.Locale)

		return nil, nil
	}

	fmt.Printf("Found diff: %s...\n", diff)      //nolint:forbidigo // need for lambda logs
	fmt.Printf("Updating %s...\n", param.Locale) //nolint:forbidigo // need for lambda logs

	// Create Feed
	feed := createRiotGamesJobsFeed(param, entries)

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
