package github

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/bradfitz/slice"
	"github.com/google/go-github/github"
	"github.com/streambinder/solbump/resource"
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

	tags, _, err := client().Repositories.ListTags(ctx, address.User, address.Project, &github.ListOptions{PerPage: 1500})
	if err != nil {
		return "", "", err
	}

	if len(tags) == 0 {
		return "", "", fmt.Errorf("No tag found")
	}

	// FIXME: as soon as GitHub APIs start ordering tags
	// by creation date, drop this block
	for _, tag := range tags {
		if tag.Commit.Author == nil {
			commit, _, err := client().Repositories.GetCommit(ctx, address.User, address.Project, tag.GetCommit().GetSHA())
			if err == nil {
				tag.Commit = commit.GetCommit()
			}
		}
	}
	slice.Sort(tags[:], func(i, j int) bool {
		return (tags[i].GetCommit().GetAuthor().GetDate().After(tags[j].GetCommit().GetAuthor().GetDate()))
	})

	return strings.ReplaceAll(url, address.Tag, tags[0].GetName()), resource.StripVersion(tags[0].GetName()), nil
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
