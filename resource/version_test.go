package resource

import "testing"

var versions = map[string]string{
	"1.0.0":        "1.0.0",
	"prefix-1.0.0": "1.0.0",
	"1.0":          "1.0",
	"prefix-1.0":   "1.0",
	"1":            "1",
	"v1":           "1",
}

func TestStripVersion(t *testing.T) {
	for key, val := range versions {
		if version := StripVersion(key); version != val {
			t.Errorf("Unexpected %s stripped version: expected %s, got %s", key, val, version)
		}
	}
}
