package gitlab

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/streambinder/solbump/resource"
	"github.com/xanzy/go-gitlab"
)

var (
	regTarball = regexp.MustCompile(`(?m)^http[s]://gitlab.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/\-/archive/(?P<Release>[a-zA-Z0-9\-\.]+)/.*.tar.gz$`)
)

// ReleaseTarballProvider represents the Provider implementation
// corresponding to gitlab.com release tarball
type ReleaseTarballProvider struct {
}

type tarballAddress struct {
	address
	Release string
}

// Name returns the name ID of the provider
func (provider ReleaseTarballProvider) Name() string {
	return "Gitlab Release tarball"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider ReleaseTarballProvider) Ready() error {
	if len(os.Getenv(GitlabEnvironmentKey)) == 0 {
		return fmt.Errorf("Gitlab API key not found")
	}

	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider ReleaseTarballProvider) Support(url, version string) bool {
	return len(regTarball.FindStringSubmatch(url)) > 1
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider ReleaseTarballProvider) Bump(url, hash, version string) (string, string, error) {
	address, err := parseTarballAddress(url)
	if err != nil {
		return "", "", err
	}

	tags, _, err := client().Tags.ListTags(address.User+"/"+address.Project, &gitlab.ListTagsOptions{})
	if err != nil {
		return "", "", err
	}

	if len(tags) == 0 {
		return "", "", fmt.Errorf("No tag found")
	}
	tagName := tags[0].Name

	return strings.ReplaceAll(url, address.Release, tagName),
		resource.StripVersion(tagName), nil
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
