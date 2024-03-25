package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

type Config struct {
	Workers         int    `koanf:"workers"`
	OutputDirectory string `koanf:"outputDirectory`
}

// Set config directory
var k = koanf.New("/")

func LoadConfig() (Config, error) {
	var conf Config

	if err := k.Load(file.Provider("/app/config.yml"), yaml.Parser()); err != nil {
		return Config{}, err
	}

	if err := k.Unmarshal("", &conf); err != nil {
		return Config{}, err
	}

	return conf, nil
}
