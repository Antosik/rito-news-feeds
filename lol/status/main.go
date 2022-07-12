package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lol"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	parameters    []statusParameters
	parametersErr error

	locales    map[string]statusLocale
	localesErr error

	domain string
)

const (
	channelsCount = 5
)

func init() {
	parameters, parametersErr = getStatusParameters()
	locales, localesErr = getStatusLocales()
	domain = os.Getenv("DOMAIN_NAME")
}

func process(
	filesChannel chan []string,
	errorsChannel chan internal.ErrorCollector,
	parameters []statusParameters,
	locales map[string]statusLocale,
	uploader *internal.S3FeedUploader,
) {
	var (
		invalidatePaths []string
		generatedFiles  []internal.FeedFile
		errorsCollector = internal.NewErrorCollector()
	)

	for _, param := range parameters {
		client := lol.StatusClient{Region: param.ID}

		for _, locale := range param.Locales {
			fpath := internal.FormatFilePath(
				filepath.Join("lol", locale, fmt.Sprintf("status.%s", param.ID)),
			)

			// Get new items
			entries, err := client.GetItems(locale)
			if err != nil {
				errorsCollector.Collect(fmt.Errorf("can't get items for %s-%s: %w", param.ID, locale, err))
				continue
			}

			// Check diff with existing data
			existingEntries, err := internal.GetExistingRawEntries[lol.StatusEntry](domain, fpath)
			if err != nil {
				errorsCollector.Collect(err)
			} else if internal.IsEqual(existingEntries, entries, compareLolStatusEntry) {
				fmt.Printf("%s-%s doesn't require update\n", param.ID, locale)
				continue
			}

			fmt.Printf("updating %s-%s...\n", param.ID, locale)

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
			rawpath := fmt.Sprintf("%s.json", fpath)

			rawfile, err := internal.GenerateRawFile(entries, rawpath)
			if err != nil {
				errorsCollector.Collect(err)
			} else {
				files = append(files, rawfile)
			}

			invalidatePaths = append(invalidatePaths, fmt.Sprintf("/%s.*", fpath))
			generatedFiles = append(generatedFiles, files...)
		}
	}

	// Upload files to S3
	if len(generatedFiles) > 0 {
		uploaderErrors := uploader.UploadFiles(generatedFiles)
		if len(uploaderErrors) > 0 {
			errorsCollector.CollectMany(uploaderErrors)
		}
	}

	errorsChannel <- *errorsCollector
	filesChannel <- invalidatePaths
}

func handler() error {
	if len(parameters) == 0 {
		return fmt.Errorf("no parameters found: %w", parametersErr)
	}

	if len(locales) == 0 {
		return fmt.Errorf("no locales found: %w", localesErr)
	}

	if domain == "" {
		return fmt.Errorf("unable to load domain name")
	}

	var (
		errorsChannel   = make(chan internal.ErrorCollector)
		filesChannel    = make(chan []string)
		filesCollector  = []string{}
		errorsCollector = internal.NewErrorCollector()
		uploader        = internal.NewS3Uploader()
		invalidator     = internal.NewCloudFrontInvalidator()
	)

	for _, chunk := range internal.SplitSliceToChunks(parameters, channelsCount) {
		go process(filesChannel, errorsChannel, chunk, locales, uploader)
	}

	for i := 0; i < channelsCount; i++ {
		errorsCollector.CollectFrom(<-errorsChannel)

		filesCollector = append(filesCollector, <-filesChannel...)
	}

	fmt.Printf("Generated files count: %d\n", len(filesCollector))

	// Invalidate CloudFront if new files were generated
	if len(filesCollector) > 0 {
		invalidationErr := invalidator.Invalidate(
			fmt.Sprintf("lolstatus-%v", time.Now().UTC().Unix()),
			filesCollector,
		)
		if invalidationErr != nil {
			errorsCollector.Collect(invalidationErr)
		}
	}

	if errorsCollector.Size() > 0 {
		fmt.Printf("%d errors occured\n", errorsCollector.Size())
		fmt.Println(errorsCollector.Error())
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
