package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StorageType    string `env:"STORAGE_TYPE" env-default:"memory"`
	DataBaseUrl    string `env:"DATABASE_URL"`
	Port           string `env:"PORT" env-default:"8080"`
	MemoryCapacity int    `env:"MEMORY_CAPACITY" env-default:"100000"`
}

func Load() (*Config, error) {
	var cfg Config
	_ = cleanenv.ReadConfig(".env", &cfg)
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("failed to read environment: %w", err)
	}
	return &cfg, nil
}
