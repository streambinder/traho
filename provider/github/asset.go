package github

import (
	"fmt"
	"regexp"

	"github.com/agnivade/levenshtein"
)

var (
	regAsset = regexp.MustCompile(`(?m)^http[s]?://github.com/(?P<User>[a-zA-Z0-9\-]+)/(?P<Project>.+)/releases/download/(?P<Release>[a-zA-Z0-9\-\.]+)/(?P<Resource>.+)$`)
)

// AssetProvider represents the Provider implementation
// corresponding to github.com release asset
type AssetProvider struct {
}

type assetAddress struct {
	address
	Release  string
	Resource string
}

// Name returns the name ID of the provider
func (provider AssetProvider) Name() string {
	return "Github Release asset"
}

// Support returns true if the given url string
// is supported by the provider
func (provider AssetProvider) Support(url, version string) bool {
	_, err := parseAssetAddress(url)
	return err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider AssetProvider) Bump(url, hash, version string) (string, string, error) {
	address, err := parseAssetAddress(url)
	if err != nil {
		return "", "", err
	}

	rels, _, err := client().Repositories.ListReleases(ctx, address.User, address.Project, nil)
	if err != nil {
		return "", "", err
	}

	if len(rels) == 0 {
		return "", "", fmt.Errorf("No release found")
	}

	var (
		bumpedURL = ""
		tagName   = address.Release
	)
	for _, rel := range rels {
		if *rel.Prerelease || *rel.Draft {
			continue
		}

		if address.Release >= *rel.TagName {
			break
		}

		for _, asset := range rel.Assets {
			if levenshtein.ComputeDistance(url, *asset.BrowserDownloadURL) <
				levenshtein.ComputeDistance(url, bumpedURL) {
				bumpedURL = *asset.BrowserDownloadURL
				tagName = *rel.TagName
			}
		}
	}

	return bumpedURL, tagName, nil
}

// Hashes returns whether or not the provider uses
// source mapping value of a source as an hash
func (provider AssetProvider) Hashes() bool {
	return true
}

func parseAssetAddress(url string) (*assetAddress, error) {
	asset := regAsset.FindStringSubmatch(url)
	if len(asset) < 5 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	address := new(assetAddress)
	address.User = asset[1]
	address.Project = asset[2]
	address.Release = asset[3]
	address.Resource = asset[4]
	return address, nil
}
