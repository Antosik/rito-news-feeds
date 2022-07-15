package internal

import (
	"bytes"
	"context"
	"fmt"
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

type S3FeedClient struct {
	s3Uploader   *manager.Uploader
	s3Downloader *manager.Downloader
}

func (client S3FeedClient) UploadFile(file FeedFile) error {
	_, err := client.s3Uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(file.Name),
		Body:        bytes.NewReader(file.Buffer),
		ContentType: aws.String(file.MimeType),
	})

	return err
}

func (client S3FeedClient) UploadFiles(files []FeedFile) []error {
	errors := make([]error, 0, len(files))

	for _, file := range files {
		err := client.UploadFile(file)
		if err != nil {
			errors = append(errors, fmt.Errorf("can't upload file %s: %w", file.Name, err))
		}
	}

	return errors
}

func (client S3FeedClient) DownloadFile(path string) (FeedFile, error) {
	buffer := manager.NewWriteAtBuffer([]byte{})

	_, err := client.s3Downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(path),
	})
	if err != nil {
		return FeedFile{}, fmt.Errorf("can't download file %s: %w", path, err)
	}

	return FeedFile{
		Name:     path,
		MimeType: "application/json",
		Buffer:   buffer.Bytes(),
	}, nil
}

func NewS3Client() *S3FeedClient {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)
	}

	if bucketName == "" {
		log.Fatalf("unable to load bucket name")
	}

	s3Client := s3.NewFromConfig(cfg)
	s3Uploader := manager.NewUploader(s3Client)
	s3Downloader := manager.NewDownloader(s3Client)

	return &S3FeedClient{s3Uploader: s3Uploader, s3Downloader: s3Downloader}
}
