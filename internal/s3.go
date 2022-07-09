package internal

import (
	"bytes"
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	bucketName string
)

func init() {
	bucketName = os.Getenv("BUCKET_NAME")
}

type S3FeedUploader struct {
	s3Manager *manager.Uploader
}

func (uploader S3FeedUploader) UploadFile(file FeedFile) error {
	_, err := uploader.s3Manager.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(file.Name),
		Body:        bytes.NewReader(file.Buffer),
		ContentType: aws.String(file.MimeType),
	})

	return err
}

func (uploader S3FeedUploader) UploadFiles(files []FeedFile) []error {
	errors := make([]error, 0, len(files))

	for _, file := range files {
		err := uploader.UploadFile(file)
		if err != nil {
			errors = append(errors, err)
		}
	}

	return errors
}

func NewS3Uploader() *S3FeedUploader {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}
	if bucketName == "" {
		log.Fatalf("unable to load bucket name")
	}

	s3Client := s3.NewFromConfig(cfg)
	s3Manager := manager.NewUploader(s3Client)

	return &S3FeedUploader{s3Manager: s3Manager}
}
