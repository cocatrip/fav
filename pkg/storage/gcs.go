package storage

type GcsBucket struct {
	Bucket string `mapstructure:"bucket"`
	Path   string `mapstructure:"path,omitempty"`
}
