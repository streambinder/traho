package fileserver

import (
	"testing"
)

func TestParseVariableAddress(t *testing.T) {
	var (
		base         = "https://davidepucci.it/traho"
		release      = "1.0.0"
		resource     = "1.0.0.tar.bz2"
		url          = base + "/" + release + "/" + resource
		address, err = parseVariableAddress(url)
	)

	if err != nil {
		t.Errorf("Unexpected failure while parsing %s", url)
	}

	if address.Base != base {
		t.Errorf("Unexpected base mismatch: expected %s, got %s", base, address.Base)
	}

	if address.Release != release {
		t.Errorf("Unexpected release mismatch: expected %s, got %s", release, address.Release)
	}

	if address.Resource != resource {
		t.Errorf("Unexpected resource mismatch: expected %s, got %s", resource, address.Resource)
	}
}
