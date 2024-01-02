package storage

type S3Bucket struct {
	Region string `mapstructure:"region"`
	Bucket string `mapstructure:"bucket"`
	Path   string `mapstructure:"path,omitempty"`
}
