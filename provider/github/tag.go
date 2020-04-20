package github

import (
	"fmt"
	"regexp"

	"github.com/Masterminds/semver"
	"github.com/bradfitz/slice"
	"github.com/google/go-github/github"
	"github.com/streambinder/solbump/resource"
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

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider TagProvider) Ready() error {
	if err := envReady(); err != nil {
		return err
	}

	return nil
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

	tags, _, err := client().Repositories.ListTags(ctx, address.User, address.Project, &github.ListOptions{PerPage: 1500})
	if err != nil {
		return "", "", err
	}

	if len(tags) == 0 {
		return "", "", fmt.Errorf("No tag found")
	}

	slice.Sort(tags[:], func(i, j int) bool {
		versionFirst, errFirst := semver.NewVersion(resource.StripVersion(tags[i].GetName()))
		if errFirst != nil {
			return false
		}
		versionLatter, errLatter := semver.NewVersion(resource.StripVersion(tags[j].GetName()))
		if errLatter != nil {
			return true
		}
		return versionFirst.GreaterThan(versionLatter)
	})

	return url, tags[0].GetName(), nil
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
	address.User = tag[1]
	address.Project = tag[2]
	return address, nil
}
