package provider

import (
	"fmt"

	"github.com/streambinder/solbump/provider/github"
	"github.com/streambinder/solbump/provider/gitlab"
)

// Provider represents a source from which a package
// and/or a tarball can be obtained
type Provider interface {
	Name() string
	Ready() error
	Support(string) bool
	Bump(string) (string, string, error)
}

var definedProviders = []Provider{
	new(github.ReleaseAssetProvider),
	new(github.ReleaseTarballProvider),
	new(github.TagTarballProvider),
	new(gitlab.ReleaseTarballProvider),
}

// All return the array of usable providers
func All() (readyProviders []Provider) {
	for _, provider := range definedProviders {
		if provider.Ready() == nil {
			readyProviders = append(readyProviders, provider)
		}
	}

	return
}

// For returns the corresponding Provider
// for the given url, an error otherwise
func For(url string) (providers []Provider, err error) {
	for _, provider := range All() {
		if provider.Support(url) {
			providers = append(providers, provider)
		}
	}

	if len(providers) == 0 {
		err = fmt.Errorf("No suitable provider for the given url")
	}

	return
}
