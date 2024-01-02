package config

import "github.com/cocatrip/fav/pkg/storage"

type Config struct {
	Storages []storage.Storage `mapstructure:"storages"`
}
