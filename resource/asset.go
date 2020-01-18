package resource

import (
	"fmt"
	"path/filepath"
)

// Asset represents a single package definition
// corresponding to a package.yml file
type Asset struct {
	// package.yml fields
	Name    string              `yaml:"name"`
	Version string              `yaml:"version"`
	Release int                 `yaml:"release"`
	Source  []map[string]string `yaml:"source"`

	// auxiliary fields
	Filename    string
	BumpVersion string
	BumpRelease int
	BumpSource  []map[string]string
}

// SourceID returns a simplified
// ID for a given source entry
func (asset *Asset) SourceID(url string) string {
	return fmt.Sprintf("%s:%s", asset.Name, filepath.Base(url))
}
