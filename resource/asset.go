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

// ID returns a nice way to represent
// a specific asset
func (asset *Asset) ID() string {
	return fmt.Sprintf("%s:%s", asset.Name, asset.Version)
}

// SourceID returns a simplified
// ID for a given source entry
func SourceID(url string) string {
	return filepath.Base(url)
}
