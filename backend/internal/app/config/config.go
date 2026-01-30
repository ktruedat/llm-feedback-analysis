package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"go.yaml.in/yaml/v3"
)

type Config struct {
	Profile     Profile     `yaml:"profile" envPrefix:"PROFILE_"`
	Server      Server      `yaml:"server" envPrefix:"SERVER_"`
	Pagination  Pagination  `yaml:"pagination" envPrefix:"PAGINATION_"`
	DB          Database    `yaml:"database" envPrefix:"DATABASE_"`
	Tracing     Tracing     `yaml:"tracing" envPrefix:"TRACING_"`
	JWT         JWT         `yaml:"jwt" envPrefix:"JWT_"`
	LLMAnalysis LLMAnalysis `yaml:"llm_analysis" envPrefix:"LLM_ANALYSIS_"`
}

func New(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			fmt.Printf("failed to close config file: %v\n", err)
		}
	}(f)

	var cfg Config
	decoder := yaml.NewDecoder(f)
	if err = decoder.Decode(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse yaml config file: %w", err)
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("could not process env vars: %w", err)
	}

	if err = cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

type validatable interface {
	Validate() error
}

func (c *Config) Validate() error {
	components := []validatable{
		c.Server,
		c.Pagination,
		c.DB,
		c.Tracing,
		c.Profile,
		c.JWT,
		c.LLMAnalysis,
	}

	for _, v := range components {
		if err := v.Validate(); err != nil {
			return err
		}
	}

	return nil
}
