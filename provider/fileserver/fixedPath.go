package fileserver

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/PuerkitoBio/goquery"
	"github.com/streambinder/solbump/resource"
)

var (
	regFixedPath = regexp.MustCompile(`^(?P<Base>http[s]?://.*)/(?P<Resource>[a-zA-Z0-9\.\-]+)$`)
)

// FixedPathProvider represents the Provider implementation
// corresponding to a generic file server asset with
// fixed parent path
type FixedPathProvider struct {
}

// Name returns the name ID of the provider
func (provider FixedPathProvider) Name() string {
	return "File Server fixed path asset"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider FixedPathProvider) Ready() error {
	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider FixedPathProvider) Support(url, version string) bool {
	address, err := parseAddress(url)
	if err != nil {
		return false
	}

	_, err = semver.NewVersion(resource.StripVersion(address.Resource))
	return err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider FixedPathProvider) Bump(url, hash, version string) (string, string, error) {
	address, err := parseAddress(url)
	if err != nil {
		return url, version, err
	}
	resourceName := strings.Split(address.Resource, version)[0]

	doc, err := goquery.NewDocument(address.Base)
	if err != nil {
		return url, version, err
	}

	var versions []href
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		url, _ := item.Attr("href")
		if len(resourceName) < len(item.Text()) &&
			resourceName == item.Text()[:len(resourceName)] &&
			filepath.Ext(url) == filepath.Ext(address.Resource) {
			versions = append(versions, href{Label: item.Text(), Href: url})
		}
	})

	if len(versions) == 0 {
		return url, version, fmt.Errorf("No relase found")
	}

	sort.SliceStable(versions, func(i, j int) bool {
		versionFirst, errFirst := semver.NewVersion(resource.StripVersion(versions[i].Href))
		if errFirst != nil {
			return false
		}
		versionLatter, errLatter := semver.NewVersion(resource.StripVersion(versions[j].Href))
		if errLatter != nil {
			return true
		}
		return versionFirst.GreaterThan(versionLatter)
	})

	return address.Base + "/" + versions[0].Href, resource.StripVersion(versions[0].Href), nil
}

// Hashes returns whether or not the provider uses
// source mapping value of a source as an hash
func (provider FixedPathProvider) Hashes() bool {
	return true
}

func parseAddress(url string) (*address, error) {
	regAsset := regFixedPath.FindStringSubmatch(url)
	if len(regAsset) < 3 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	addressAsset := new(address)
	addressAsset.Base = regAsset[1]
	addressAsset.Resource = regAsset[2]
	return addressAsset, nil
}
