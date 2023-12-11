package config

import (
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Topics map[string][]string `yaml:"topics"`
}

var configFolder = "./config"

func ReadConfigFile() (*Config, error) {
	files, err := os.ReadDir(configFolder)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, nil
	}
	if len(files) > 1 {
		return nil, errors.New("expected only 1 config file but found multiple")
	}
	configData, err := os.ReadFile(path.Join(configFolder, files[0].Name()))
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
