package gitlab

import (
	"context"
	"fmt"

	"github.com/streambinder/traho/config"
	"github.com/xanzy/go-gitlab"
)

type address struct {
	Proto, Domain, User, Project string
}

var (
	cli *gitlab.Client
	ctx = context.Background()
)

func client(proto, domain string) *gitlab.Client {
	if cli != nil {
		return cli
	}
	cli, _ = gitlab.NewClient(
		config.Get().Gitlab.API,
		gitlab.WithBaseURL(fmt.Sprintf("%s://%s", proto, domain)))
	return cli
}
