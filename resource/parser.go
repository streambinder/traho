package resource

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Parse generates a new Asset instance
// starting from a package.yml filename
func Parse(fname string) (*Asset, error) {
	asset := new(Asset)

	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&asset); err != nil {
		return nil, err
	}

	asset.Filename = fname
	return asset, nil
}
