package provider

import (
	"testing"

	"github.com/streambinder/solbump/provider/fileserver"
	"github.com/streambinder/solbump/provider/github"
	"github.com/streambinder/solbump/provider/gitlab"
	"github.com/streambinder/solbump/provider/undated"
)

const version = "1.0.0"

var sources = map[string]string{
	"git|https://github.com/streambinder/solbump.git":                                                  new(github.TagProvider).Name(),
	"https://github.com/streambinder/solbump/releases/download/" + version + "/asset.zip":              new(github.AssetProvider).Name(),
	"https://github.com/streambinder/solbump/archive/" + version + ".tar.gz":                           new(github.TarballProvider).Name(),
	"https://gitlab.com/streambinder/solbump/-/archive/" + version + "/solbump-" + version + ".tar.gz": new(gitlab.Provider).Name(),
	"https://davidepucci.it/solbump/" + version + "/solbump-" + version + ".tar.bz2":                   new(fileserver.VariablePathProvider).Name(),
	"https://davidepucci.it/solbump/solbump-" + version + ".tar.bz2":                                   new(fileserver.FixedPathProvider).Name(),
	"https://davidepucci.it/download?project=Solbump":                                                  new(undated.Provider).Name(),
}

func TestFor(t *testing.T) {
	for key, val := range sources {
		prov, err := For(key, version)
		if err != nil {
			t.Errorf("Unexpected failure while getting provider for %s: %s", key, err)
		}

		if prov.Name() != val {
			t.Errorf("Unexpected %s provider: expected %s, got %s", key, val, prov.Name())
		}
	}
}
