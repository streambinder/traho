package github

import (
	"testing"
)

func TestParseTarballAddress(t *testing.T) {
	var (
		user         = "streambinder"
		project      = "solbump"
		tag          = "1.0.0"
		url          = "https://github.com/" + user + "/" + project + "/archive/" + tag + ".tar.gz"
		tarball, err = parseTarballAddress(url)
	)

	if err != nil {
		t.Errorf("Unexpected failure while parsing %s", url)
	}

	if tarball.User != user {
		t.Errorf("Unexpected user mismatch: expected %s, got %s", user, tarball.User)
	}

	if tarball.Project != project {
		t.Errorf("Unexpected project mismatch: expected %s, got %s", project, tarball.Project)
	}

	if tarball.Tag != tag {
		t.Errorf("Unexpected tag mismatch: expected %s, got %s", tag, tarball.Tag)
	}
}
