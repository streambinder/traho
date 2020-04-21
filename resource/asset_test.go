package resource

import (
	"testing"
)

func TestID(t *testing.T) {
	var (
		name    = "TestID"
		version = "v1"
		asset   = &Asset{Name: name, Version: version}
	)
	if asset.ID() != name+":"+version {
		t.Errorf("Unexpected asset ID")
	}
}

func TestSourceID(t *testing.T) {
	if SourceID("https://davidepucci.it/index.html") != "index.html" {
		t.Errorf("Unexpected source ID")
	}
}
