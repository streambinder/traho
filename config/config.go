package config

// Config represents the abstraction of the parsed
// configuration file
type Config struct {
	Github *Github `yaml:"github"`
}

// Github represents the configuration section
// related to the Github provider
type Github struct {
	API string `yaml:"api"`
}
