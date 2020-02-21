package github

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var (
	regTagTarball = regexp.MustCompile(`(?m)^http[s]://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/archive/(?P<Tag>[a-zA-Z0-9\-\.]+).tar.gz$`)
)

// TagTarballProvider represents the Provider implementation
// corresponding to github.com tag tarball
type TagTarballProvider struct {
}

type tagTarballAddress struct {
	address
	Tag string
}

// Name returns the name ID of the provider
func (provider TagTarballProvider) Name() string {
	return "Github Tag tarball"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider TagTarballProvider) Ready() error {
	if len(os.Getenv(GithubEnvironmentKey)) != 40 {
		return fmt.Errorf("Github API key not found")
	}

	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider TagTarballProvider) Support(url, version string) bool {
	return len(regTarball.FindStringSubmatch(url)) > 1
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider TagTarballProvider) Bump(url string) (string, string, error) {
	address, err := parseTagTarballAddress(url)
	if err != nil {
		return "", "", err
	}

	tags, _, err := client().Repositories.ListTags(ctx, address.User, address.Project, nil)
	if err != nil {
		return "", "", err
	}

	var tagName string
	for _, tag := range tags {
		if address.Tag >= *tag.Name {
			break
		}

		tagName = *tag.Name
		break
	}

	if len(tagName) == 0 {
		return "", "", fmt.Errorf("No release found")
	}

	return strings.ReplaceAll(url, address.Tag, tagName),
		regVersionStrip.ReplaceAllString(tagName, ""), nil
}

func parseTagTarballAddress(url string) (*tagTarballAddress, error) {
	regTarball := regTagTarball.FindStringSubmatch(url)
	if len(regTarball) < 4 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	addressTarball := new(tagTarballAddress)
	addressTarball.User = regTarball[1]
	addressTarball.Project = regTarball[2]
	addressTarball.Tag = regTarball[3]
	return addressTarball, nil
}
