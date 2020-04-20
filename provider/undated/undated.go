package undated

import (
	"strings"
	"time"

	"github.com/streambinder/solbump/resource"
)

const (
	versionFormat = "20060102"
)

// Provider represents the Provider implementation
// corresponding to a generic undated asset
type Provider struct {
}

// Name returns the name ID of the provider
func (provider Provider) Name() string {
	return "Undated resource"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider Provider) Ready() error {
	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider Provider) Support(url, version string) bool {
	_, err := time.Parse(versionFormat, version)
	return !strings.Contains(url, version) && strings.Contains(url, "?") && err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider Provider) Bump(url, hash, version string) (string, string, error) {
	if bumpHash, err := resource.Hash(url); err == nil && bumpHash == hash {
		return url, version, nil
	}

	return url, time.Now().Format(versionFormat), nil
}

// Hashes returns whether or not the provider uses
// source mapping value of a source as an hash
func (provider Provider) Hashes() bool {
	return true
}
