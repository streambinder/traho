package provider

import (
	"fmt"

	"github.com/streambinder/solbump/provider/fileserver"
	"github.com/streambinder/solbump/provider/github"
	"github.com/streambinder/solbump/provider/gitlab"
	"github.com/streambinder/solbump/provider/undated"
)

// Provider represents a source from which a package
// and/or a tarball can be obtained
type Provider interface {
	Name() string
	Support(string, string) bool
	Bump(string, string, string) (string, string, error)
	Hashes() bool
}

var definedProviders = []Provider{
	new(github.TagProvider),
	new(github.AssetProvider),
	new(github.TarballProvider),
	new(gitlab.Provider),
	new(fileserver.VariablePathProvider),
	new(fileserver.FixedPathProvider),
	new(undated.Provider),
}

// All return the array of usable providers
func All() (readyProviders []Provider) {
	return definedProviders
}

// For returns the corresponding Provider
// for the given url, an error otherwise
func For(url, version string) (Provider, error) {
	for _, provider := range All() {
		if provider.Support(url, version) {
			return provider, nil
		}
	}

	return nil, fmt.Errorf("No provider for given url")
}
