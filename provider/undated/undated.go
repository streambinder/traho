package undated

import (
	"strings"
	"time"
)

const (
	versionFormat = "20060102"
)

// ResourceProvider represents the Provider implementation
// corresponding to a generic undated asset
type ResourceProvider struct {
}

// Name returns the name ID of the provider
func (provider ResourceProvider) Name() string {
	return "Undated resource"
}

// Ready returns an error if the provider
// is unconfigured or unusable
func (provider ResourceProvider) Ready() error {
	return nil
}

// Support returns true if the given url string
// is supported by the provider
func (provider ResourceProvider) Support(url, version string) bool {
	_, err := time.Parse(versionFormat, version)
	return !strings.Contains(url, version) && strings.Contains(url, "?") && err == nil
}

// Bump returns the bump of the given url and
// the updated associated version or, if unable, an error
func (provider ResourceProvider) Bump(url string) (string, string, error) {
	return url, time.Now().Format(versionFormat), nil
}
