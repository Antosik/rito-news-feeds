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

func init() {
	parameters, parametersErr = getStatusParameters()
	locales, localesErr = getStatusLocales()

	domain = os.Getenv("DOMAIN_NAME")
}

func process(
	channel chan internal.ErrorCollector,
	parameters []statusParameters,
	locales map[string]statusLocale,
	uploader *internal.S3FeedUploader,
) {
	var (
		generatedFiles  []internal.FeedFile
		errorsCollector = internal.NewErrorCollector()
	)

	for _, param := range parameters {
		var (
			client = lol.StatusClient{Region: param.Id}
			dpath  = filepath.Join("lol", param.Region)
		)

		for _, locale := range param.Locales {
			entries, err := client.GetItems(locale)
			if err != nil {
				errorsCollector.Collect(fmt.Errorf("can't get items for %s-%s: %w", param.Id, locale, err))
				continue
			}

			fpath := filepath.Join(dpath, fmt.Sprintf("status.%s", locale))

			localeData, ok := locales[locale]
			if !ok {
				localeData = locales["en_US"]
			}

			// Create Feed
			feed := createLolStatusFeed(param.Id, &localeData, entries)

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
	}

	uploaderErrors := uploader.UploadFiles(generatedFiles)
	if len(uploaderErrors) > 0 {
		errorsCollector.CollectMany(uploaderErrors)
	}

	channel <- *errorsCollector
}

func handler() error {
	if len(parameters) == 0 {
		return fmt.Errorf("no parameters found: %w", parametersErr)
	}
	if len(locales) == 0 {
		return fmt.Errorf("no locales found: %w", localesErr)
	}

	var (
		channel         = make(chan internal.ErrorCollector)
		channelsCount   = 5
		errorsCollector = internal.NewErrorCollector()
		uploader        = internal.NewS3Uploader()
		invalidator     = internal.NewCloudFrontInvalidator()
	)

	for _, chunk := range internal.SplitSliceToChunks(parameters, channelsCount) {
		go process(channel, chunk, locales, uploader)
	}

	for i := 0; i < channelsCount; i++ {
		errorsCollector.CollectFrom(<-channel)
	}

	invalidationErr := invalidator.Invalidate(
		fmt.Sprintf("lolstatus-%v", time.Now().UTC().Unix()),
		[]string{"/lol/*/status*"},
	)
	if invalidationErr != nil {
		errorsCollector.Collect(invalidationErr)
	}

	if errorsCollector.Size() > 0 {
		fmt.Printf(errorsCollector.Error())
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
