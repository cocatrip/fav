package storage

import (
	"os"
)

type Storage struct {
	Name      string    `mapstructure:"name"`
	Type      string    `mapstructure:"type"`
	S3Bucket  S3Bucket  `mapstructure:"s3,omitempty"`
	GcsBucket GcsBucket `mapstructure:"gcs,omitempty"`
}

func (s *Storage) getBucketType() {
	if s.S3Bucket != (S3Bucket{}) {
		s.Type = "s3"
		return
	}

	if s.GcsBucket != (GcsBucket{}) {
		s.Type = "gcs"
		return
	}
}

func (s *Storage) Upload(name string) error {
	s.getBucketType()

	switch s.Type {
	case "gcs":
		if err := s.GcsBucket.Upload(name); err != nil {
			return err
		}
	case "s3":
	}

	os.Remove(name)

	return nil
}

func (s *Storage) Download(name string) error {
	s.getBucketType()

	switch s.Type {
	case "gcs":
		s.GcsBucket.Download(name)
	case "s3":
	}

	return nil
}
