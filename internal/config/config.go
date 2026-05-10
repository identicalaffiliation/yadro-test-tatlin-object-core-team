package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type CLIConfig struct {
	ChannelBuff int `yaml:"chan_buffer"`
	BatchSize   int `yaml:"batch_size"`
	BucketCount int `yaml:"bucket_count"`
}

func LoadCLIConfig(configPath string) (*CLIConfig, error) {
	CLIConfig := new(CLIConfig)
	if err := cleanenv.ReadConfig(configPath, CLIConfig); err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	return CLIConfig, nil
}
