package provider

import (
	"testing"

	"github.com/streambinder/solbump/provider/fileserver"
	"github.com/streambinder/solbump/provider/github"
	"github.com/streambinder/solbump/provider/gitlab"
	"github.com/streambinder/solbump/provider/undated"
)

type source struct {
	url      string
	version  string
	provider string
}

const version = "1.0.0"

var sources = []source{
	{"git|https://github.com/streambinder/solbump.git", version, new(github.TagProvider).Name()},
	{"https://github.com/streambinder/solbump/releases/download/" + version + "/asset.zip", version, new(github.AssetProvider).Name()},
	{"https://github.com/streambinder/solbump/archive/" + version + ".tar.gz", version, new(github.TarballProvider).Name()},
	{"https://gitlab.com/streambinder/solbump/-/archive/" + version + "/solbump-" + version + ".tar.gz", version, new(gitlab.Provider).Name()},
	{"https://davidepucci.it/solbump/" + version + "/solbump-" + version + ".tar.bz2", version, new(fileserver.VariablePathProvider).Name()},
	{"https://davidepucci.it/solbump/solbump-" + version + ".tar.bz2", version, new(fileserver.FixedPathProvider).Name()},
	{"https://davidepucci.it/download?project=Solbump", "20060102", new(undated.Provider).Name()},
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
