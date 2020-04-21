package config

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

// Parse generates a new Config instance
// starting from a configuration file path
func Parse(fname string) error {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		config = new(Config)
		return nil
	}

	content, err := ioutil.ReadFile(fname)
	if err != nil {
		return err
	}

	return parseYaml(content)
}

func parseYaml(content []byte) error {
	config = new(Config)
	if err := yaml.Unmarshal(content, config); err != nil {
		return err
	}

	return nil
}
