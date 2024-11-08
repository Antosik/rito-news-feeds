package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lor"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	params    []statusParameters
	errParams error

	locales    map[string]statusLocale
	errLocales error

	domain string

	statusProcessor LoRStatusProcessor
	mainProcessor   internal.MainProcessor[statusParameters]
)

const (
	channelsCount = 10
)

func init() {
	params, errParams = getStatusParameters()
	locales, errLocales = getStatusLocales()
	domain = os.Getenv("DOMAIN_NAME")

	statusProcessor = LoRStatusProcessor{}
	mainProcessor = internal.MainProcessor[statusParameters]{
		Name:        "lorstatus",
		Concurrency: channelsCount,

		TypeProcessor: &statusProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// LoR Status Processor (implements AbstractProcessor).
type LoRStatusProcessor struct{}

func (p *LoRStatusProcessor) GenerateFilePath(param statusParameters, locale string) string {
	return internal.FormatFilePath(filepath.Join("lor", locale, "status."+param.ID))
}

func (p *LoRStatusProcessor) GenerateAbstractFilePath(param statusParameters) string {
	return filepath.Join("/", "lor", "*", "status."+param.ID+".*")
}

func (p *LoRStatusProcessor) ProcessParameters(
	param statusParameters,
) ([]internal.FeedFile, []error) {
	var (
		generatedFiles  []internal.FeedFile
		client          = lor.StatusClient{Region: param.ID}
		errorsCollector = internal.NewErrorCollector()
	)

	for _, locale := range locales {
		fpath := p.GenerateFilePath(param, locale.Locale)

		// Get new items
		entries, err := client.GetItems(locale.Locale)
		if err != nil {
			errorsCollector.Collect(fmt.Errorf("can't get items for %s-%s: %w", param.ID, locale.Locale, err))

			continue
		}

		// Check diff with existing data
		rawpath := fpath + ".json"

		existingFile, err := mainProcessor.S3Client.DownloadFile(rawpath)
		if err != nil {
			errorsCollector.Collect(err)
		}

		existingEntries, err := internal.GetExistingRawEntries[lor.StatusEntry](existingFile)
		if err != nil {
			errorsCollector.Collect(err)
		}

		diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getLorStatusEntryKey)
		if isEqual {
			//nolint:forbidigo // need for lambda logs
			fmt.Printf("%s-%s doesn't require update\n", param.ID, locale.Locale)

			return nil, nil
		}

		fmt.Printf("Found diff: %s...\n", diff)                    //nolint:forbidigo // need for lambda logs
		fmt.Printf("Updating %s-%s...\n", param.ID, locale.Locale) //nolint:forbidigo // need for lambda logs

		// Create Feed
		feed := createLorStatusFeed(param.ID, locale, entries)

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

		generatedFiles = append(generatedFiles, files...)
	}

	return generatedFiles, *errorsCollector
}

func handler() error {
	if len(params) == 0 {
		return fmt.Errorf("no params found: %w", errParams)
	}

	if len(locales) == 0 {
		return fmt.Errorf("no locales found: %w", errLocales)
	}

	if domain == "" {
		return internal.ErrDomainNotFound
	}

	return mainProcessor.Process(params)
}

func main() {
	lambda.Start(handler)
}
