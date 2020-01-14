package github

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	regTarball = regexp.MustCompile(`(?m)^http[s]://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/archive/(?P<Release>[a-zA-Z0-9\-\.]+).tar.gz$`)
)

// ReleaseTarballProvider represents the Provider implementation
// corresponding to github.com release tarball
type ReleaseTarballProvider struct {
}

type tarballAddress struct {
	address
	Release string
}

// Name returns the name ID of the provider
func (provider ReleaseTarballProvider) Name() string {
	return "Github Release tarball"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider ReleaseTarballProvider) Ready() error {
	if len(os.Getenv(GithubEnvironmentKey)) != 40 {
		return fmt.Errorf("Github API key not found")
	}

	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider ReleaseTarballProvider) Support(url string) bool {
	return len(regTarball.FindStringSubmatch(url)) > 1
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider ReleaseTarballProvider) Bump(url string) (string, string, error) {
	address, err := parseTarballAddress(url)
	if err != nil {
		return "", "", err
	}

	rels, _, err := client().Repositories.ListReleases(ctx, address.User, address.Project, nil)
	if err != nil {
		return "", "", err
	}

	var tagName string
	for _, rel := range rels {
		if *rel.Prerelease || *rel.Draft {
			continue
		}
		tagName = *rel.TagName
	}

	if len(tagName) == 0 {
		return "", "", fmt.Errorf("No release found")
	}

	return strings.ReplaceAll(url, address.Release, tagName),
		regVersionStrip.ReplaceAllString(tagName, ""), nil
}

func parseTarballAddress(url string) (*tarballAddress, error) {
	regTarball := regTarball.FindStringSubmatch(url)
	if len(regTarball) < 4 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	addressTarball := new(tarballAddress)
	addressTarball.User = regTarball[1]
	addressTarball.Project = regTarball[2]
	addressTarball.Release = regTarball[3]
	return addressTarball, nil
}
