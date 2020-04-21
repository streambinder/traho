package github

import (
	"context"
	"fmt"
	"regexp"

	"github.com/google/go-github/github"
	"github.com/streambinder/solbump/config"
	"golang.org/x/oauth2"
)

var (
	regAddress = regexp.MustCompile(`(?m)^http[s]://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/.*$`)
)

type address struct {
	User, Project string
}

var (
	cli *github.Client
	ctx = context.Background()
)

func envReady() error {
	if len(config.Get().Github.API) != 40 {
		return fmt.Errorf("Github API key not found")
	}

	return nil
}

func parseAddress(url string) (*address, error) {
	regURL := regAddress.FindStringSubmatch(url)
	if len(regURL) < 3 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	return &address{User: regURL[1], Project: regURL[2]}, nil
}

func client() *github.Client {
	if cli != nil {
		return cli
	}

	cli = github.NewClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Get().Github.API},
	)))
	return cli
}
