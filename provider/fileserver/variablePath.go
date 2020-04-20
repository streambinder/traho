package fileserver

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
	"github.com/PuerkitoBio/goquery"
	"github.com/bradfitz/slice"
	"github.com/streambinder/solbump/resource"
)

var (
	regVariablePath = regexp.MustCompile(`^(?P<Base>http[s]?://.*)/(?P<Release>[0-9]{1}[0-9\.]+)/(?P<Resource>[a-zA-Z0-9\.\-]+)$`)
)

// VariablePathProvider represents the Provider implementation
// corresponding to a generic file server asset with
// variable parent path
type VariablePathProvider struct {
}

type variableAddress struct {
	address
	Release string
}

// Name returns the name ID of the provider
func (provider VariablePathProvider) Name() string {
	return "File Server variable path asset"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider VariablePathProvider) Ready() error {
	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider VariablePathProvider) Support(url, version string) bool {
	_, err := parseVariableAddress(url)
	return err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider VariablePathProvider) Bump(url, hash, version string) (string, string, error) {
	address, err := parseVariableAddress(url)
	if err != nil {
		return url, version, err
	}

	doc, err := goquery.NewDocument(address.Base)
	if err != nil {
		return url, version, err
	}

	var paths []href
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		url, _ := item.Attr("href")
		paths = append(paths, href{Label: item.Text(), Href: url})
	})

	if len(paths) == 0 {
		return url, version, fmt.Errorf("No relase found")
	}

	slice.Sort(paths[:], func(i, j int) bool {
		versionFirst, errFirst := semver.NewVersion(resource.StripVersion(paths[i].Label))
		if errFirst != nil {
			return false
		}
		versionLatter, errLatter := semver.NewVersion(resource.StripVersion(paths[j].Label))
		if errLatter != nil {
			return true
		}
		return versionFirst.GreaterThan(versionLatter)
	})

	doc, err = goquery.NewDocument(address.Base + "/" + paths[0].Href)
	if err != nil {
		return url, version, err
	}

	var versions []href
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		url, _ := item.Attr("href")
		if filepath.Ext(url) == filepath.Ext(address.Resource) {
			versions = append(versions, href{Label: item.Text(), Href: url})
		}
	})

	slice.Sort(versions[:], func(i, j int) bool {
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

	return address.Base + "/" + strings.ReplaceAll(paths[0].Href+"/"+versions[0].Href, "//", "/"),
		resource.StripVersion(versions[0].Href), nil
}

// Hashes returns whether or not the provider uses
// source mapping value of a source as an hash
func (provider VariablePathProvider) Hashes() bool {
	return true
}

func parseVariableAddress(url string) (*variableAddress, error) {
	regAsset := regVariablePath.FindStringSubmatch(url)
	if len(regAsset) < 4 {
		return nil, fmt.Errorf("Unrecognized url %s", url)
	}

	addressAsset := new(variableAddress)
	addressAsset.Base = regAsset[1]
	addressAsset.Release = regAsset[2]
	addressAsset.Resource = regAsset[3]
	return addressAsset, nil
}
