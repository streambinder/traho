package resource

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"

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

// Flush updates the physical content of a given asset
// with bumped fields
func Flush(asset *Asset) error {
	read, err := ioutil.ReadFile(asset.Filename)
	if err != nil {
		return err
	}

	var (
		content    = string(read)
		regVersion = regexp.MustCompile(`(?m)(^version\s+:)\s+` + regexp.QuoteMeta(asset.Version) + `$`)
		regRelease = regexp.MustCompile(`(?m)(^release\s+:)\s+` + regexp.QuoteMeta(strconv.Itoa(asset.Release)) + `$`)
	)
	content = regVersion.ReplaceAllString(content, `$1 `+asset.BumpVersion)
	content = regRelease.ReplaceAllString(content, `$1 `+strconv.Itoa(asset.BumpRelease))
	for entry := range asset.Source {
		for source, hash := range asset.Source[entry] {
			for bumpSource, bumpHash := range asset.BumpSource[entry] {
				fmt.Println(regexp.QuoteMeta(source), regexp.QuoteMeta(hash))
				fmt.Println(bumpSource, bumpHash)
				regSource := regexp.MustCompile(`(\s{4}\-)\s+` + regexp.QuoteMeta(source) + `\s+:\s+` + regexp.QuoteMeta(hash))
				content = regSource.ReplaceAllString(content, `$1 `+bumpSource+` : `+bumpHash)
			}
		}
	}

	err = ioutil.WriteFile(asset.Filename, []byte(content), 0)
	if err != nil {
		return err
	}

	return nil
}
