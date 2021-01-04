package config

import (
	"strings"
)

// Config represents the abstraction of the parsed
// configuration file
type Config struct {
	Github struct {
		API string `yaml:"api"`
	} `yaml:"github"`
	Gitlab struct {
		API     string `yaml:"api"`
		Domains []struct {
			Name string `yaml:"name"`
			API  string `yaml:"api"`
		} `yaml:"domains"`
	} `yaml:"gitlab"`
}

var config *Config

// Get current config instance
func Get() *Config {
	return config
}

// HasGitlabDomain returns true if config has the given
// domain defined as custom Gitlab domain
func HasGitlabDomain(domain string) bool {
	domain = strings.ToLower(domain)
	for _, customDomain := range config.Gitlab.Domains {
		if strings.ToLower(customDomain.Name) == domain {
			return true
		}
	}
	return false
}
