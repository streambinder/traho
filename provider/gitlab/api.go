package gitlab

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"github.com/xanzy/go-gitlab"
)

// GitlabEnvironmentKey is the environment key associated to
// Gitlab API key
const GitlabEnvironmentKey = "solbump.provider.gitlab"

var (
	regAddress = regexp.MustCompile(`(?m)^http[s]://gitlab.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/.*$`)
)

type address struct {
	User, Project string
}

var (
	cli *gitlab.Client
	ctx = context.Background()
)

func parseAddress(url string) (*address, error) {
	regURL := regAddress.FindStringSubmatch(url)
	if len(regURL) < 3 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	return &address{User: regURL[1], Project: regURL[2]}, nil
}

func client() *gitlab.Client {
	if cli != nil {
		return cli
	}

	cli, _ = gitlab.NewClient(os.Getenv(GitlabEnvironmentKey))
	return cli
}
