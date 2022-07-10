package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/lol"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	parameters    []newsParameters
	parametersErr error

	domain string
)

const (
	articlesCount = 100
	channelsCount = 5
)

func init() {
	parameters, parametersErr = getNewsParameters()
	domain = os.Getenv("DOMAIN_NAME")
}

func process(
	filesChannel chan int,
	errorsChannel chan internal.ErrorCollector,
	parameters []newsParameters,
	uploader *internal.S3FeedUploader,
) {
	var (
		generatedFiles  []internal.FeedFile
		errorsCollector = internal.NewErrorCollector()
	)

	for _, param := range parameters {
		var (
			client = lol.NewsClient{Locale: strings.ToLower(param.Locale)}
			dpath  = filepath.Join("lol", param.Region)
		)

		fpath := internal.FormatFilePath(filepath.Join(dpath, fmt.Sprintf("news.%s", param.Locale)))

		// Get new items
		entries, err := client.GetItems(articlesCount)
		if err != nil {
			errorsCollector.Collect(fmt.Errorf("can't get items for %s-%s: %w", param.Region, param.Locale, err))
			continue
		}

		// Check diff with existing data
		existingEntries, err := internal.GetExistingRawEntries[lol.NewsEntry](domain, fpath)
		if err != nil {
			errorsCollector.Collect(err)
		} else if internal.IsEqual(existingEntries, entries, compareLolNewsEntry) {
			fmt.Printf("%s-%s doesn't require update\n", param.Region, param.Locale)
			continue
		}

		fmt.Printf("updating %s-%s...\n", param.Region, param.Locale)

		// Create Feed
		feed := createLolNewsFeed(param, entries)

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

	// Upload files to S3
	if len(generatedFiles) > 0 {
		uploaderErrors := uploader.UploadFiles(generatedFiles)
		if len(uploaderErrors) > 0 {
			errorsCollector.CollectMany(uploaderErrors)
		}
	}

	errorsChannel <- *errorsCollector
	filesChannel <- len(generatedFiles)
}

func handler() error {
	if len(parameters) == 0 {
		return fmt.Errorf("no parameters found: %w", parametersErr)
	}

	if domain == "" {
		return fmt.Errorf("unable to load domain name")
	}

	var (
		errorsChannel       = make(chan internal.ErrorCollector)
		filesChannel        = make(chan int)
		generatedFilesCount = 0
		errorsCollector     = internal.NewErrorCollector()
		uploader            = internal.NewS3Uploader()
		invalidator         = internal.NewCloudFrontInvalidator()
	)

	for _, chunk := range internal.SplitSliceToChunks(parameters, channelsCount) {
		go process(filesChannel, errorsChannel, chunk, uploader)
	}

	for i := 0; i < channelsCount; i++ {
		errorsCollector.CollectFrom(<-errorsChannel)

		generatedFilesCount = generatedFilesCount + <-filesChannel
	}

	fmt.Printf("Generated files count: %d\n", generatedFilesCount)

	// Invalidate CloudFront if new files were generated
	if generatedFilesCount > 0 {
		invalidationErr := invalidator.Invalidate(
			fmt.Sprintf("lolnews-%v", time.Now().UTC().Unix()),
			[]string{"/lol/*/news*"},
		)
		if invalidationErr != nil {
			errorsCollector.Collect(invalidationErr)
		}
	}

	if errorsCollector.Size() > 0 {
		fmt.Printf("%d errors occured\n", errorsCollector.Size())
		fmt.Printf(errorsCollector.Error())
	}

	return nil
}

func main() {
	lambda.Start(handler)
}
