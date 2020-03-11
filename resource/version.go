package resource

import "regexp"

var (
	regVersion = regexp.MustCompile("[^0-9\\.]+")
)

// StripVersion returns the normalized variant string
// corresponding to the given version string
func StripVersion(version string) string {
	return regVersion.ReplaceAllString(version, "")
}
