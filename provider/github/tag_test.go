package github

import (
	"testing"
)

func TestParseTagAddress(t *testing.T) {
	var (
		user     = "streambinder"
		project  = "solbump"
		url      = "git|https://github.com/" + user + "/" + project + ".git"
		tag, err = parseTagAddress(url)
	)

	if err != nil {
		t.Errorf("Unexpected failure while parsing %s", url)
	}

	if tag.User != user {
		t.Errorf("Unexpected user mismatch: expected %s, got %s", user, tag.User)
	}

	if tag.Project != project {
		t.Errorf("Unexpected project mismatch: expected %s, got %s", project, tag.Project)
	}
}
