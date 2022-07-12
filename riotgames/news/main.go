package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Antosik/rito-news-feeds/internal"
	"github.com/Antosik/rito-news/riotgames"

	"github.com/aws/aws-lambda-go/lambda"
)

var (
	parameters    []newsParameters
	parametersErr error

	domain string
)

const (
	articlesCount = 100
	channelsCount = 1
)

func init() {
	parameters, parametersErr = getNewsParameters()
	domain = os.Getenv("DOMAIN_NAME")
}

func process(
	filesChannel chan []string,
	errorsChannel chan internal.ErrorCollector,
	parameters []newsParameters,
	uploader *internal.S3FeedUploader,
) {
	var (
		invalidatePaths []string
		generatedFiles  []internal.FeedFile
		errorsCollector = internal.NewErrorCollector()
	)

	for _, param := range parameters {
		fmt.Printf("start processing %s\n", param.Locale)

		var (
			client = riotgames.NewsClient{Locale: strings.ToLower(param.Locale)}
			fpath  = internal.FormatFilePath(filepath.Join("riotgames", param.Locale, "news"))
		)

		// Get new items
		entries, err := client.GetItems(articlesCount)
		if err != nil {
			errorsCollector.Collect(fmt.Errorf("can't get items for %s: %w", param.Locale, err))
			continue
		}

		// Check diff with existing data
		existingEntries, err := internal.GetExistingRawEntries[riotgames.NewsEntry](domain, fpath)
		if err != nil {
			errorsCollector.Collect(err)
		} else if internal.IsEqual(existingEntries, entries, compareRiotGamesNewsEntry) {
			fmt.Printf("%s doesn't require update\n", param.Locale)
			continue
		}

		fmt.Printf("updating %s...\n", param.Locale)

		// Create Feed
		feed := createRiotGamesNewsFeed(param, entries)

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

	// Upload files to S3
	if len(generatedFiles) > 0 {
		uploaderErrors := uploader.UploadFiles(generatedFiles)
		if len(uploaderErrors) > 0 {
			errorsCollector.CollectMany(uploaderErrors)
		}
	}

	if len(invalidatePaths) > len(parameters)/3 {
		invalidatePaths = []string{
			internal.FormatFilePath(filepath.Join("/", "riotgames", "*", "news.*")),
		}
	}

	errorsChannel <- *errorsCollector
	filesChannel <- invalidatePaths
}

func handler() error {
	if len(parameters) == 0 {
		return fmt.Errorf("no parameters found: %w", parametersErr)
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
		go process(filesChannel, errorsChannel, chunk, uploader)
	}

	for i := 0; i < channelsCount; i++ {
		errorsCollector.CollectFrom(<-errorsChannel)

		filesCollector = append(filesCollector, <-filesChannel...)
	}

	fmt.Printf("Generated files count: %d\n", len(filesCollector))

	// Invalidate CloudFront if new files were generated
	if len(filesCollector) > 0 {
		invalidationErr := invalidator.Invalidate(
			fmt.Sprintf("riotgamesnews-%v", time.Now().UTC().Unix()),
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
