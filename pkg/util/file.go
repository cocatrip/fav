package util

import (
	"os"
)

func ReadFile(name string) (*string, error) {
	buf, err := os.ReadFile(name)

	if err != nil {
		return nil, err
	}

	result := string(buf)

	return &result, nil
}
