package main

import (
	"fmt"
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
)

func init() {
	parameters, parametersErr = getStatusParameters()
	locales, localesErr = getStatusLocales()
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
		client := lol.StatusClient{Region: param.Id}

		dpath := filepath.Join("lol", param.Region)
		for _, locale := range param.Locales {
			entries, err := client.GetItems(locale)
			if err != nil {
				errorsCollector.Collect(fmt.Errorf("can't get items for %s-%s: %w", param.Id, locale, err))
				continue
			}

			feed := createLolStatusFeed(locale, entries)

			fpath := filepath.Join(dpath, fmt.Sprintf("status.%s", locale))
			files, errors := internal.GenerateFeedFiles(feed, fpath)
			if len(errors) > 0 {
				errorsCollector.CollectMany(errors)
			}

			rawpath := filepath.Join(dpath, fmt.Sprintf("status.%s.json", locale))
			rawjson, err := internal.MarshalJSON(entries)
			if err != nil {
				errorsCollector.Collect(err)
			} else {
				files = append(files, internal.FeedFile{
					Name:     rawpath,
					MimeType: "application/json",
					Buffer:   rawjson,
				})
			}

			generatedFiles = append(generatedFiles, files...)
		}

		time.Sleep(500 * time.Millisecond)
	}

	uploaderErrors := uploader.UploadFiles(generatedFiles)
	if len(uploaderErrors) > 0 {
		errorsCollector.CollectMany(uploaderErrors)
	}

	channel <- *errorsCollector
}

func StatusHandler() error {
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
	)

	for _, chunk := range internal.SplitSliceToChunks(parameters, channelsCount) {
		go process(channel, chunk, locales, uploader)
	}

	for i := 0; i < channelsCount; i++ {
		errorsCollector.CollectFrom(<-channel)
	}

	if errorsCollector.Size() > 0 {
		fmt.Printf(errorsCollector.Error())
	}

	return nil
}

func main() {
	lambda.Start(StatusHandler)
}
