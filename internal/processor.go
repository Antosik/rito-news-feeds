package internal

import (
	"fmt"
)

type AbstractProcessor[ParamType interface{}] interface {
	GenerateAbstractFilePath(parameters ParamType) string
	ProcessParameters(parameters ParamType) ([]FeedFile, []error)
}

type MainProcessor[ParamType interface{}] struct {
	// Main props
	Name        string
	Concurrency int

	// DI
	TypeProcessor AbstractProcessor[ParamType]

	// AWS
	S3Client *S3FeedClient
}

func (p *MainProcessor[ParamType]) ProcessChunk(
	filesChannel chan []string,
	errorsChannel chan ErrorCollector,
	parametersChunk []ParamType,
	s3Client *S3FeedClient,
) {
	var (
		generatedPaths  []string
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

			generatedPaths = append(
				generatedPaths,
				p.TypeProcessor.GenerateAbstractFilePath(param),
			)
		}
	}

	// Upload files to S3
	if len(generatedFiles) > 0 {
		s3ClientErrors := s3Client.UploadFiles(generatedFiles)
		if len(s3ClientErrors) > 0 {
			errorsCollector.CollectMany(s3ClientErrors)
		}
	}

	errorsChannel <- *errorsCollector
	filesChannel <- generatedPaths
}

func (p *MainProcessor[ParamType]) Process(params []ParamType) error {
	var (
		errorsChannel           = make(chan ErrorCollector)
		errorsCollector         = NewErrorCollector()
		generatedPathsChannel   = make(chan []string)
		generatedPathsCollector = []string{}
	)

	fmt.Printf(
		"Starting processing %s with %d parameters and %d concurrency\n",
		p.Name,
		len(params),
		p.Concurrency,
	)

	for _, chunk := range SplitSliceToChunks(params, p.Concurrency) {
		go p.ProcessChunk(generatedPathsChannel, errorsChannel, chunk, p.S3Client)
	}

	for i := 0; i < p.Concurrency; i++ {
		errorsCollector.CollectFrom(<-errorsChannel)

		generatedPathsCollector = append(generatedPathsCollector, <-generatedPathsChannel...)
	}

	if len(generatedPathsCollector) > 0 {
		fmt.Printf("Updated paths (%d): %v\n", len(generatedPathsCollector), generatedPathsCollector)
	}

	if errorsCollector.Size() > 0 {
		fmt.Printf("During execution %d errors occured\n", errorsCollector.Size())
		return errorsCollector.Error()
	}

	return nil
}
