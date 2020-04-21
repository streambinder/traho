package gitlab

import (
	"testing"
)

func TestParseTarballAddress(t *testing.T) {
	var (
		user         = "streambinder"
		project      = "solbump"
		release      = "1.0.0"
		url          = "https://gitlab.com/" + user + "/" + project + "/-/archive/" + release + "/" + project + "-" + release + ".tar.gz"
		tarball, err = parseTarballAddress(url)
	)

	if err != nil {
		t.Errorf("Unexpected failure while parsing %s", url)
	}

	if tarball.User != user {
		t.Errorf("Unexpected user mismatch: expected %s, got %s", user, tarball.User)
	}

	if tarball.Project != project {
		t.Errorf("Unexpected user mismatch: expected %s, got %s", project, tarball.Project)
	}

	if tarball.Release != release {
		t.Errorf("Unexpected user mismatch: expected %s, got %s", release, tarball.Release)
	}
}
