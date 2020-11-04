package fileserver

import (
	"testing"
)

func TestParseAddress(t *testing.T) {
	var (
		base         = "https://davidepucci.it/traho"
		resource     = "1.0.0.tar.bz2"
		url          = base + "/" + resource
		address, err = parseAddress(url)
	)

	if err != nil {
		t.Errorf("Unexpected failure while parsing %s", url)
	}

	if address.Base != base {
		t.Errorf("Unexpected base mismatch: expected %s, got %s", base, address.Base)
	}

	if address.Resource != resource {
		t.Errorf("Unexpected resource mismatch: expected %s, got %s", resource, address.Resource)
	}
}
