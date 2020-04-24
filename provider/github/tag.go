package github

import (
	"fmt"
	"regexp"
)

var (
	regTag = regexp.MustCompile(`(?m)^git\|http[s]://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+).git$`)
)

// TagProvider represents the Provider implementation
// corresponding to github.com tag tarball
type TagProvider struct {
}

type tagAddress struct {
	address
	Tag string
}

// Name returns the name ID of the provider
func (provider TagProvider) Name() string {
	return "Github tag"
}

// Support returns true if the given url string
// is supported by the provider
func (provider TagProvider) Support(url, version string) bool {
	_, err := parseTagAddress(url)
	return err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider TagProvider) Bump(url, tag, version string) (string, string, error) {
	address, err := parseTagAddress(url)
	if err != nil {
		return "", "", err
	}
	address.Tag = tag

	release, err := getLatestRelease(address.User, address.Project)
	if err == nil {
		return address.Full, *release.TagName, nil
	}

	repoTag, err := getLatestTag(address.User, address.Project)
	if err != nil {
		return "", "", err
	}

	return url, repoTag.GetName(), nil
}

// Hashes returns whether or not the provider uses
// source mapping value of a source as an hash
func (provider TagProvider) Hashes() bool {
	return false
}

func parseTagAddress(url string) (*tagAddress, error) {
	tag := regTag.FindStringSubmatch(url)
	if len(tag) != 3 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	address := new(tagAddress)
	address.Full = url
	address.User = tag[1]
	address.Project = tag[2]
	return address, nil
}
