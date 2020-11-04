package provider

import (
	"testing"

	"github.com/streambinder/traho/provider/fileserver"
	"github.com/streambinder/traho/provider/github"
	"github.com/streambinder/traho/provider/gitlab"
	"github.com/streambinder/traho/provider/undated"
)

type source struct {
	url      string
	version  string
	provider string
}

const version = "1.0.0"

var sources = []source{
	{"git|https://github.com/streambinder/traho.git", version, new(github.TagProvider).Name()},
	{"https://github.com/streambinder/traho/releases/download/" + version + "/asset.zip", version, new(github.AssetProvider).Name()},
	{"https://github.com/streambinder/traho/archive/" + version + ".tar.gz", version, new(github.TarballProvider).Name()},
	{"https://gitlab.com/streambinder/traho/-/archive/" + version + "/traho-" + version + ".tar.gz", version, new(gitlab.Provider).Name()},
	{"https://davidepucci.it/traho/" + version + "/traho-" + version + ".tar.bz2", version, new(fileserver.VariablePathProvider).Name()},
	{"https://davidepucci.it/traho/traho-" + version + ".tar.bz2", version, new(fileserver.FixedPathProvider).Name()},
	{"https://davidepucci.it/download?project=traho", "20060102", new(undated.Provider).Name()},
}

func TestFor(t *testing.T) {
	for _, source := range sources {
		prov, err := For(source.url, source.version)
		if err != nil {
			t.Errorf("Unexpected failure while getting provider for %s, version, %s", source.url, err)
		}

		if prov.Name() != source.provider {
			t.Errorf("Unexpected %s provider, version, expected %s, got %s", source.url, source.provider, prov.Name())
		}
	}
}
