package gitlab

import (
	"context"
	"fmt"
	"regexp"

	"github.com/streambinder/traho/config"
	"github.com/xanzy/go-gitlab"
)

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

	cli, _ = gitlab.NewClient(config.Get().Gitlab.API)
	return cli
}
