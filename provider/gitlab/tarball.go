package gitlab

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/streambinder/traho/config"
	"github.com/streambinder/traho/resource"
	"github.com/xanzy/go-gitlab"
)

var (
	regTarball = regexp.MustCompile(`(?m)^(?P<Proto>http[s]?)://(?P<Domain>[a-zA-Z0-9\.\-]+)/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/\-/archive/(?P<Release>[a-zA-Z0-9\-\.]+)/.*.(zip|tar.gz|tar.bz2|tar)$`)
)

// Provider represents the Provider implementation
// corresponding to gitlab.com release tarball
type Provider struct {
}

type tarballAddress struct {
	address
	Release string
}

// Name returns the name ID of the provider
func (provider Provider) Name() string {
	return "Gitlab Release tarball"
}

// Support returns true if the given url string
// is supported by the provider
func (provider Provider) Support(url, version string) bool {
	_, err := parseTarballAddress(url)
	return err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider Provider) Bump(url, hash, version string) (string, string, error) {
	address, err := parseTarballAddress(url)
	if err != nil {
		return "", "", err
	}

	tags, _, err := client(address.Proto, address.Domain).Tags.ListTags(address.User+"/"+address.Project, &gitlab.ListTagsOptions{})
	if err != nil {
		logrus.Println(err)
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
	if len(regTarball) < 6 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	addr := &tarballAddress{address{Proto: regTarball[1], Domain: regTarball[2], User: regTarball[3], Project: regTarball[4]}, regTarball[5]}
	if !(addr.Domain == "gitlab.com" || config.HasGitlabDomain(addr.Domain)) {
		return nil, fmt.Errorf("Domain %s does not point to a Gitlab instance", addr.Domain)
	}

	return addr, nil
}

// Hashes returns whether or not the provider uses
// source mapping value of a source as an hash
func (provider Provider) Hashes() bool {
	return true
}
