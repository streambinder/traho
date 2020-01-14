package config

import (
	"os"

	"github.com/streambinder/solbump/provider/github"
	"gopkg.in/yaml.v2"
)

// Parse generates a new Config instance
// starting from a configuration file path
func Parse(fname string) (*Config, error) {
	config := new(Config)
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return config, nil
	}

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, err
	}

	return config, config.process()
}

func (config *Config) process() error {
	if len(config.Github.API) > 0 {
		if err := os.Setenv(
			github.GithubEnvironmentKey,
			config.Github.API,
		); err != nil {
			return err
		}
	}

	return nil
}
