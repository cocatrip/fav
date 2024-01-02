package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	awsS3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type S3Bucket struct {
	Region string `mapstructure:"region"`
	Bucket string `mapstructure:"bucket"`
	Path   string `mapstructure:"path,omitempty"`
}

func (s3 *S3Bucket) Upload(object string) error {
	file, err := os.Open(object)
	if err != nil {
		return err
	}
	defer file.Close()

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3.Region)},
	)

	if _, err := sess.Config.Credentials.Get(); err != nil {
		return err
	}

	if len(s3.Path) != 0 {
		object = filepath.Join(s3.Path, object)
	}

	uploader := s3manager.NewUploader(sess)

	uploadInput := &s3manager.UploadInput{
		Bucket: aws.String(s3.Bucket),
		Key:    aws.String(object),
		Body:   file,
	}

	if _, err := uploader.Upload(uploadInput); err != nil {
		return err
	}

	return nil
}

func (s3 *S3Bucket) Download(prefix string) error {
	var object string

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(s3.Region)},
	)
	if err != nil {
		return err
	}

	svc := awsS3.New(sess)

	resp, err := svc.ListObjectsV2(&awsS3.ListObjectsV2Input{Bucket: aws.String(s3.Bucket)})
	if err != nil {
		fmt.Printf("Unable to list items in bucket %q, %v", s3.Bucket, err)
	}

	if len(s3.Path) != 0 {
		prefix = filepath.Join(s3.Path, prefix)
	}

	for _, item := range resp.Contents {
		if strings.HasPrefix(*item.Key, prefix) {
			object = *item.Key
		}
	}

	destFileName := filepath.Join(".tmp", object)

	if err := os.MkdirAll(filepath.Dir(destFileName), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(destFileName)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer f.Close()

	downloader := s3manager.NewDownloader(sess)

	getObjectInput := &awsS3.GetObjectInput{
		Bucket: aws.String(s3.Bucket),
		Key:    aws.String(object),
	}

	if _, err := downloader.Download(f, getObjectInput); err != nil {
		return err
	}

	deleteObjectInput := &awsS3.DeleteObjectInput{
		Bucket: aws.String(s3.Bucket),
		Key:    aws.String(object),
	}
	if _, err = svc.DeleteObject(deleteObjectInput); err != nil {
		return err
	}

	return nil
}
