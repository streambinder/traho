package config

// Config represents the abstraction of the parsed
// configuration file
type Config struct {
	Github struct {
		API string `yaml:"api"`
	} `yaml:"github"`
	Gitlab struct {
		API string `yaml:"api"`
	} `yaml:"gitlab"`
}
