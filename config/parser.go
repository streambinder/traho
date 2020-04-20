package config

import (
	"io/ioutil"
	"os"

	"github.com/streambinder/solbump/provider/github"
	"github.com/streambinder/solbump/provider/gitlab"
	"gopkg.in/yaml.v2"
)

// Parse generates a new Config instance
// starting from a configuration file path
func Parse(fname string) (*Config, error) {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return new(Config), nil
	}

	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return nil, err
	}

	return parseYaml(content)
}

func parseYaml(content []byte) (*Config, error) {
	config := new(Config)
	if err := yaml.Unmarshal(content, config); err != nil {
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

	if len(config.Gitlab.API) > 0 {
		if err := os.Setenv(
			gitlab.GitlabEnvironmentKey,
			config.Gitlab.API,
		); err != nil {
			return err
		}
	}

	return nil
}
