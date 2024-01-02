package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type GcsBucket struct {
	Bucket string `mapstructure:"bucket"`
	Path   string `mapstructure:"path,omitempty"`
}

func (gcs *GcsBucket) Upload(object string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %s", err)
	}
	defer client.Close()

	f, err := os.Open(object)
	if err != nil {
		return fmt.Errorf("os.Open: %s", err)
	}
	defer f.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	if len(gcs.Path) != 0 {
		object = filepath.Join(gcs.Path, object)
	}

	o := client.Bucket(gcs.Bucket).Object(object)

	o = o.If(storage.Conditions{DoesNotExist: true})

	wc := o.NewWriter(ctx)
	if _, err = io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (gcs *GcsBucket) Download(prefix string) error {
	var object string

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %w", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if len(gcs.Path) != 0 {
		prefix = filepath.Join(gcs.Path, prefix)
	}

	it := client.Bucket(gcs.Bucket).Objects(ctx, &storage.Query{
		Prefix:    prefix,
		Delimiter: "_",
	})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Bucket(%q).Objects(): %w", gcs.Bucket, err)
		}

		object = attrs.Name
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

	rc, err := client.Bucket(gcs.Bucket).Object(object).NewReader(ctx)
	if err != nil {
		return fmt.Errorf("Object(%q).NewReader: %w", object, err)
	}
	defer rc.Close()

	if _, err := io.Copy(f, rc); err != nil {
		return fmt.Errorf("io.Copy: %w", err)
	}

	if err = f.Close(); err != nil {
		return fmt.Errorf("f.Close: %w", err)
	}

	o := client.Bucket(gcs.Bucket).Object(object)

	attrs, err := o.Attrs(ctx)
	if err != nil {
		return fmt.Errorf("object.Attrs: %w", err)
	}
	o = o.If(storage.Conditions{GenerationMatch: attrs.Generation})

	if err := o.Delete(ctx); err != nil {
		return fmt.Errorf("Object(%q).Delete: %w", object, err)
	}

	return nil
}
