package github

import (
	"fmt"
	"os"
	"regexp"

	"github.com/agnivade/levenshtein"
)

var (
	regAsset = regexp.MustCompile(`(?m)^http[s]://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/releases/download/(?P<Release>[a-zA-Z0-9\-\.]+)/(?P<Resource>.+)$`)
)

// ReleaseAssetProvider represents the Provider implementation
// corresponding to github.com release asset
type ReleaseAssetProvider struct {
}

type assetAddress struct {
	address
	Release  string
	Resource string
}

// Name returns the name ID of the provider
func (provider ReleaseAssetProvider) Name() string {
	return "Github Release asset"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider ReleaseAssetProvider) Ready() error {
	if len(os.Getenv(GithubEnvironmentKey)) != 40 {
		return fmt.Errorf("Github API key not found")
	}

	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider ReleaseAssetProvider) Support(url string) bool {
	return len(regAsset.FindStringSubmatch(url)) > 1
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider ReleaseAssetProvider) Bump(url string) (string, string, error) {
	address, err := parseAssetAddress(url)
	if err != nil {
		return "", "", err
	}

	rels, _, err := client().Repositories.ListReleases(ctx, address.User, address.Project, nil)
	if err != nil {
		return "", "", err
	}

	var bumpedURL string
	for _, asset := range rels[0].Assets {
		if levenshtein.ComputeDistance(url, *asset.BrowserDownloadURL) <
			levenshtein.ComputeDistance(url, bumpedURL) {
			bumpedURL = *asset.BrowserDownloadURL
		}
	}

	if len(bumpedURL) == 0 {
		return "", "", fmt.Errorf("Unable to find a new asset")
	}

	return bumpedURL, *rels[0].TagName, nil
}

func parseAssetAddress(url string) (*assetAddress, error) {
	regAsset := regAsset.FindStringSubmatch(url)
	if len(regAsset) < 5 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	addressAsset := new(assetAddress)
	addressAsset.User = regAsset[1]
	addressAsset.Project = regAsset[2]
	addressAsset.Release = regAsset[3]
	addressAsset.Resource = regAsset[4]
	return addressAsset, nil
}
