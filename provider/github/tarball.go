package github

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	regTarball = regexp.MustCompile(`(?m)^http[s]://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/archive/(?P<Tag>[a-zA-Z0-9\-\.]+).tar.gz$`)
)

// TarballProvider represents the Provider implementation
// corresponding to github.com tag tarball
type TarballProvider struct {
}

type tagTarballAddress struct {
	address
	Tag string
}

// Name returns the name ID of the provider
func (provider TarballProvider) Name() string {
	return "Github tarball"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider TarballProvider) Ready() error {
	if len(os.Getenv(GithubEnvironmentKey)) != 40 {
		return fmt.Errorf("Github API key not found")
	}

	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider TarballProvider) Support(url, version string) bool {
	return len(regTarball.FindStringSubmatch(url)) > 1
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider TarballProvider) Bump(url, hash, version string) (string, string, error) {
	address, err := parseTagTarballAddress(url)
	if err != nil {
		return "", "", err
	}

	tags, _, err := client().Repositories.ListTags(ctx, address.User, address.Project, nil)
	if err != nil {
		return "", "", err
	}

	if len(tags) == 0 {
		return "", "", fmt.Errorf("No tag found")
	}
	tagName := *tags[0].Name

	return strings.ReplaceAll(url, address.Tag, tagName),
		regVersionStrip.ReplaceAllString(tagName, ""), nil
}

func parseTagTarballAddress(url string) (*tagTarballAddress, error) {
	regTarball := regTarball.FindStringSubmatch(url)
	if len(regTarball) < 4 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	addressTarball := new(tagTarballAddress)
	addressTarball.User = regTarball[1]
	addressTarball.Project = regTarball[2]
	addressTarball.Tag = regTarball[3]
	return addressTarball, nil
}
