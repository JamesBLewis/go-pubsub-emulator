package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectID string              `yaml:"projectID"`
	Topics    map[string][]string `yaml:"topics"`
}

func ReadConfigFile(configFile string) (map[string][]string, error) {
	configData, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		return nil, err
	}

	return config.Topics, nil
}
