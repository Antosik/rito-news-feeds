package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lol"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	params    []statusParameters
	paramsErr error

	locales    map[string]statusLocale
	localesErr error

	domain string

	statusProcessor LoLStatusProcessor
	mainProcessor   internal.MainProcessor[statusParameters]
)

const (
	channelsCount = 5
)

func init() {
	params, paramsErr = getStatusParameters()
	locales, localesErr = getStatusLocales()
	domain = os.Getenv("DOMAIN_NAME")

	statusProcessor = LoLStatusProcessor{}
	mainProcessor = internal.MainProcessor[statusParameters]{
		Name:        "lolstatus",
		Concurrency: channelsCount,

		TypeProcessor: &statusProcessor,

		S3Client: internal.NewS3Client(),
	}
}

// LoL Status Processor (implements AbstractProcessor)
type LoLStatusProcessor struct{}

func (p *LoLStatusProcessor) GenerateFilePath(param statusParameters, locale string) string {
	return internal.FormatFilePath(filepath.Join("lol", locale, fmt.Sprintf("status.%s", param.Region)))
}

func (p *LoLStatusProcessor) GenerateAbstractFilePath(param statusParameters) string {
	return filepath.Join("/", "lol", "*", fmt.Sprintf("status.%s.*", param.Region))
}

func (p *LoLStatusProcessor) ProcessParameters(
	param statusParameters,
) ([]internal.FeedFile, []error) {
	var (
		generatedFiles  []internal.FeedFile
		client          = lol.StatusClient{Region: param.ID}
		errorsCollector = internal.NewErrorCollector()
	)

	for _, locale := range param.Locales {
		fpath := p.GenerateFilePath(param, locale)

		// Get new items
		entries, err := client.GetItems(locale)
		if err != nil {
			errorsCollector.Collect(fmt.Errorf("can't get items for %s-%s: %w", param.ID, locale, err))
			continue
		}

		// Check diff with existing data
		rawpath := fmt.Sprintf("%s.json", fpath)

		existingFile, err := mainProcessor.S3Client.DownloadFile(rawpath)
		if err != nil {
			errorsCollector.Collect(err)
		}

		existingEntries, err := internal.GetExistingRawEntries[lol.StatusEntry](existingFile)
		if err != nil {
			errorsCollector.Collect(err)
		}

		diff, isEqual := internal.CompareAndGetDiff(existingEntries, entries, getLolStatusEntryKey)
		if isEqual {
			fmt.Printf("%s-%s doesn't require update\n", param.ID, locale)
			return nil, nil
		}

		fmt.Printf("Found diff: %s...\n", diff)
		fmt.Printf("Updating %s-%s...\n", param.ID, locale)

		localeData, ok := locales[locale]
		if !ok {
			localeData = locales["en_US"]
		}

		// Create Feed
		feed := createLolStatusFeed(param.ID, &localeData, entries)

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
