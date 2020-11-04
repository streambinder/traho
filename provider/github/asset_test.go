package github

import (
	"testing"
)

func TestParseAssetAddress(t *testing.T) {
	var (
		user       = "streambinder"
		project    = "traho"
		release    = "1.0.0"
		resource   = "asset.zip"
		url        = "https://github.com/" + user + "/" + project + "/releases/download/" + release + "/" + resource
		asset, err = parseAssetAddress(url)
	)

	if err != nil {
		t.Errorf("Unexpected failure while parsing %s", url)
	}

	if asset.Full != url {
		t.Errorf("Unexpected url mismatch: expected %s, got %s", url, asset.Full)
	}

	if asset.User != user {
		t.Errorf("Unexpected user mismatch: expected %s, got %s", user, asset.User)
	}

	if asset.Project != project {
		t.Errorf("Unexpected project mismatch: expected %s, got %s", project, asset.Project)
	}

	if asset.Release != release {
		t.Errorf("Unexpected release mismatch: expected %s, got %s", release, asset.Release)
	}

	if asset.Resource != resource {
		t.Errorf("Unexpected resource mismatch: expected %s, got %s", resource, asset.Resource)
	}
}
