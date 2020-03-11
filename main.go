package main

import (
	"flag"
	"os/user"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/streambinder/solbump/config"
	"github.com/streambinder/solbump/provider"
	"github.com/streambinder/solbump/resource"
	logrusPrefix "github.com/x-cray/logrus-prefixed-formatter"
)

var (
	cfg        *config.Config
	log        *logrus.Logger
	argDry     bool
	argVerbose bool
)

func init() {
	// logger setup
	log = logrus.New()
	log.Formatter = &logrusPrefix.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
		FullTimestamp:   true,
		TimestampFormat: "Mon 2 Jan 2006 15:04:05",
	}

	// config parse
	usr, err := user.Current()
	if err != nil {
		log.WithError(err).Fatalln("Unable to get current user")
	}
	if cfg, err = config.Parse(filepath.Join(usr.HomeDir, ".config/solbump")); err != nil {
		log.WithError(err).Fatalln("Unable to parse configuration file")
	}

	// args parse
	flag.BoolVar(&argDry, "dry", false, "Dry run")
	flag.BoolVar(&argVerbose, "verbose", false, "Verbose mode")
	flag.Parse()

	if argVerbose {
		log.SetLevel(logrus.DebugLevel)
	}

	if len(flag.Args()) == 0 {
		log.Println("At least a package.yml file must be given")
	}
}

func main() {
	// iterate over packages
	for _, fname := range flag.Args() {
		log.WithField("file", fname).Debugln("Inspecting resource")
		asset, err := resource.Parse(fname)
		if err != nil {
			log.WithField("file", fname).WithError(err).Warnln("Unable to parse asset")
			continue
		}

		if err := handleSources(asset); err != nil {
			log.WithField("package", asset.ID()).Errorln(err)
			continue
		}

		// check asset has been updated
		if len(asset.BumpVersion) == 0 {
			log.WithField("package", asset.ID()).Println("Package is up-to-date")
			continue
		}

		// persist update on package file
		asset.BumpRelease = asset.Release + 1

		if !argDry {
			if err := resource.Flush(asset); err != nil {
				log.WithField("file", fname).WithError(err).Errorln("Unable to flush file content")
			}
		}
		log.WithFields(logrus.Fields{
			"package": asset.ID(),
			"version": asset.BumpVersion,
		}).Println("Package updated")
	}
}

func handleSources(asset *resource.Asset) error {
	// iterate over package sources
	for idx := len(asset.Source) - 1; idx >= 0; idx-- {
		entry := asset.Source[idx]
		for url, hash := range entry {
			// pick right provider for source
			prov, err := provider.For(url, asset.Version)
			if err != nil {
				if idx == 0 {
					return err
				}
				log.WithField("asset", resource.SourceID(url)).Errorln(err)
				continue
			}

			// check source on provider for updates
			log.WithFields(logrus.Fields{
				"provider": prov.Name(),
				"asset":    resource.SourceID(url),
			}).Debugln("Analyzing asset")
			bump, version, err := prov.Bump(url, hash, asset.Version)
			if err != nil {
				if idx == 0 {
					return err
				}
				log.WithField("asset", resource.SourceID(url)).Errorln(err)
				continue
			}

			// use only first source to update version
			if asset.Version == version {
				asset.BumpSource = append(asset.BumpSource, map[string]string{url: hash})
				continue
			}

			// fetch hash sum of updated source
			if !argDry {
				log.WithField("asset", resource.SourceID(url)).Debugln("Going to fetch source and calculate SHA265 hash")
				assetHash, err := resource.Hash(bump)
				if err != nil {
					if idx == 0 {
						return err
					}
					log.WithField("asset", resource.SourceID(url)).Errorln(err)
					continue
				}
				asset.BumpSource = append(asset.BumpSource, map[string]string{bump: assetHash})
			}

			// set new package version based
			// on first source update met
			if idx == 0 {
				asset.BumpVersion = version
			}

			log.WithFields(logrus.Fields{
				"package": asset.ID(),
				"asset":   resource.SourceID(url),
				"version": version,
			}).Debugln("Source correctly bumped")
		}
	}

	return nil
}
