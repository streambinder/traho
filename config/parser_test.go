package config

import (
	"testing"
)

func TestParseEmpty(t *testing.T) {
	if _, err := parseYaml([]byte("")); err != nil {
		t.Errorf("Unexpected failure while parsing empty configuration file")
	}
}

func TestParseUnexisting(t *testing.T) {
	if _, err := Parse("/tmp/not-existing-file"); err != nil {
		t.Errorf("Unexpected failure while parsing unexisting configuration file")
	}
}
