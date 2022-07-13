package internal

import (
	"fmt"
	"time"
)

type AbstractProcessor[ParamType interface{}] interface {
	GenerateInvalidationFilePath(parameters ParamType) string
	GenerateAsteriskInvalidationPath() string
	ProcessParameters(parameters ParamType) ([]FeedFile, []error)
}

type MainProcessor[ParamType interface{}] struct {
	// Main props
	Name        string
	Concurrency int

	// DI
	TypeProcessor AbstractProcessor[ParamType]

	// AWS
	CFInvalidator *CloudFrontInvalidator
	S3Uploader    *S3FeedUploader
}

func (p *MainProcessor[ParamType]) ProcessChunk(
	filesChannel chan []string,
	errorsChannel chan ErrorCollector,
	parametersChunk []ParamType,
	uploader *S3FeedUploader,
) {
	var (
		invalidatePaths []string
		generatedFiles  []FeedFile
		errorsCollector = NewErrorCollector()
	)

	for _, param := range parametersChunk {
		files, errors := p.TypeProcessor.ProcessParameters(param)

		if len(errors) > 0 {
			errorsCollector.CollectMany(errors)
		}

		if len(files) > 0 {
			generatedFiles = append(generatedFiles, files...)

			invalidatePaths = append(
				invalidatePaths,
				p.TypeProcessor.GenerateInvalidationFilePath(param),
			)
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

func (p *MainProcessor[ParamType]) Process(params []ParamType) error {
	var (
		errorsChannel         = make(chan ErrorCollector)
		errorsCollector       = NewErrorCollector()
		invalidationChannel   = make(chan []string)
		invalidationCollector = []string{}
	)

	fmt.Printf(
		"Starting processing %s with %d parameters and %d concurrency\n",
		p.Name,
		len(params),
		p.Concurrency,
	)

	for _, chunk := range SplitSliceToChunks(params, p.Concurrency) {
		go p.ProcessChunk(invalidationChannel, errorsChannel, chunk, p.S3Uploader)
	}

	for i := 0; i < p.Concurrency; i++ {
		errorsCollector.CollectFrom(<-errorsChannel)

		invalidationCollector = append(invalidationCollector, <-invalidationChannel...)
	}

	if len(invalidationCollector) > len(params)/3 {
		invalidationCollector = []string{
			p.TypeProcessor.GenerateAsteriskInvalidationPath(),
		}
	}

	fmt.Printf("Paths to invalidate count: %d\n", len(invalidationCollector))

	// Invalidate CloudFront if new files were generated
	if len(invalidationCollector) > 0 {
		invalidationErr := p.CFInvalidator.Invalidate(
			fmt.Sprintf("%s-%v", p.Name, time.Now().UTC().Unix()),
			invalidationCollector,
		)
		if invalidationErr != nil {
			errorsCollector.Collect(invalidationErr)
		}
	}

	if errorsCollector.Size() > 0 {
		fmt.Printf("During execution %d errors occured\n", errorsCollector.Size())
		return errorsCollector.Error()
	}

	return nil
}
