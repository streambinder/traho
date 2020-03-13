package resource

import (
	"regexp"
)

var (
	regVersions = []*regexp.Regexp{
		regexp.MustCompile(`(?m)([0-9]{1}[0-9\.]+[0-9a-zA-Z\-]+)`),
		regexp.MustCompile(`(?m)^v([0-9]{1}[0-9\.]+$)`),
	}
)

// StripVersion returns the normalized variant string
// corresponding to the given version string
func StripVersion(version string) string {
	for _, reg := range regVersions {
		for _, matches := range reg.FindAllStringSubmatch(version, -1) {
			if len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return version
}
