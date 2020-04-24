package github

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	regTarball = regexp.MustCompile(`(?m)^http[s]?://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/archive/(?P<Tag>[a-zA-Z0-9\-\.]+).tar.gz$`)
)

// TarballProvider represents the Provider implementation
// corresponding to github.com tag tarball
type TarballProvider struct {
}

type tarballAddress struct {
	address
	Tag string
}

// Name returns the name ID of the provider
func (provider TarballProvider) Name() string {
	return "Github tarball"
}

// Support returns true if the given url string
// is supported by the provider
func (provider TarballProvider) Support(url, version string) bool {
	_, err := parseTarballAddress(url)
	return err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider TarballProvider) Bump(url, hash, version string) (string, string, error) {
	address, err := parseTarballAddress(url)
	if err != nil {
		return "", "", err
	}

	release, err := getLatestRelease(address.User, address.Project)
	if err == nil {
		return strings.ReplaceAll(url, address.Tag, *release.TagName), *release.TagName, nil
	}

	tag, err := getLatestTag(address.User, address.Project)
	if err != nil {
		return "", "", err
	}

	return strings.ReplaceAll(url, address.Tag, tag.GetName()), tag.GetName(), nil
}

// Hashes returns whether or not the provider uses
// source mapping value of a source as an hash
func (provider TarballProvider) Hashes() bool {
	return true
}

func parseTarballAddress(url string) (*tarballAddress, error) {
	tarball := regTarball.FindStringSubmatch(url)
	if len(tarball) < 4 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	address := new(tarballAddress)
	address.Full = url
	address.User = tarball[1]
	address.Project = tarball[2]
	address.Tag = tarball[3]
	return address, nil
}
