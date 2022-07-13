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
	paramsErr error

	locales    map[string]statusLocale
	localesErr error

	domain string

	statusProcessor LoRStatusProcessor
	mainProcessor   internal.MainProcessor[statusParameters]
)

const (
	channelsCount = 5
)

func init() {
	params, paramsErr = getStatusParameters()
	locales, localesErr = getStatusLocales()
	domain = os.Getenv("DOMAIN_NAME")

	statusProcessor = LoRStatusProcessor{}
	mainProcessor = internal.MainProcessor[statusParameters]{
		Name:        "lorstatus",
		Concurrency: channelsCount,

		TypeProcessor: &statusProcessor,

		CFInvalidator: internal.NewCloudFrontInvalidator(),
		S3Uploader:    internal.NewS3Uploader(),
	}
}

// LoR Status Processor (implements AbstractProcessor)
type LoRStatusProcessor struct{}

func (p *LoRStatusProcessor) GenerateFilePath(param statusParameters, locale string) string {
	return internal.FormatFilePath(filepath.Join("lor", locale, fmt.Sprintf("status.%s", param.ID)))
}

func (p *LoRStatusProcessor) GenerateInvalidationFilePath(param statusParameters) string {
	return filepath.Join("/", "lor", "*", fmt.Sprintf("status.%s.*", param.ID))
}

func (p *LoRStatusProcessor) GenerateAsteriskInvalidationPath() string {
	return filepath.Join("/", "lor", "*", "status.*")
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
		existingEntries, err := internal.GetExistingRawEntries[lor.StatusEntry](domain, fpath)
		if err != nil {
			errorsCollector.Collect(err)
		} else if internal.IsEqual(existingEntries, entries, compareLorStatusEntry) {
			fmt.Printf("%s-%s doesn't require update\n", param.ID, locale.Locale)
			continue
		}

		fmt.Printf("Updating %s-%s...\n", param.ID, locale.Locale)

		// Create Feed
		feed := createLorStatusFeed(param.ID, locale, entries)

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

		generatedFiles = append(generatedFiles, files...)
	}

	return generatedFiles, *errorsCollector
}

func handler() error {
	if len(params) == 0 {
		return fmt.Errorf("no params found: %w", paramsErr)
	}

	if len(locales) == 0 {
		return fmt.Errorf("no locales found: %w", localesErr)
	}

	if domain == "" {
		return internal.ErrDomainNotFound
	}

	return mainProcessor.Process(params)
}

func main() {
	lambda.Start(handler)
}
