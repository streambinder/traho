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
	cfg *config.Config
	log *logrus.Logger
	dry bool
)

func init() {
	log = logrus.New()
	log.Formatter = &logrusPrefix.TextFormatter{
		ForceColors:     true,
		ForceFormatting: true,
		FullTimestamp:   true,
		TimestampFormat: "Mon 2 Jan 2006 15:04:05",
	}

	usr, err := user.Current()
	if err != nil {
		log.WithError(err).Fatalln("Unable to get current user")
	}

	if cfg, err = config.Parse(filepath.Join(usr.HomeDir, ".config/solbump")); err != nil {
		log.WithError(err).Fatalln("Unable to parse configuration file")
	}

	flag.BoolVar(&dry, "dry", false, "Dry run")
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Println("At least a package.yml file must be given")
	}
}

func main() {
	for _, fname := range flag.Args() {
		log.WithField("file", fname).Println("Inspecting resource")
		asset, err := resource.Parse(fname)
		if err != nil {
			log.WithField("file", fname).WithError(err).Warnln("Unable to parse asset")
			continue
		}

		for entry := range asset.Source {
			for source, hash := range asset.Source[entry] {
				log.WithField("url", source).Println("Analyzing asset")

				prov, err := provider.For(source)
				if err != nil {
					log.WithField("source", source).WithError(err).Warnln("Unable to find matching provider")
					continue
				}

				log.WithFields(logrus.Fields{
					"provider": prov.Name(),
					"source":   source,
				}).Println("Bumping source")
				bump, version, err := prov.Bump(source)
				if err != nil {
					log.WithField("source", source).WithError(err).Warnln("Unable to bump source")
					continue
				}

				// use only first source to update version
				if asset.Version == version {
					asset.BumpSource = append(asset.BumpSource, map[string]string{source: hash})
					continue
				}

				log.WithField("source", bump).Println("Going to fetch source and calculate SHA265 hash")
				assetHash, err := resource.Hash(bump)
				if err != nil {
					log.WithField("source", bump).WithError(err).Warnln("Unable to calculate hash for source")
				}
				asset.BumpSource = append(asset.BumpSource, map[string]string{bump: assetHash})

				if entry == 0 {
					asset.BumpVersion = version
				}

				log.WithFields(logrus.Fields{
					"source":  bump,
					"version": version,
				}).Println("Source correctly bumped")
			}
		}

		// check asset has been updated
		if len(asset.BumpVersion) == 0 {
			log.WithField("file", fname).Println("No update detected")
			continue
		}

		if dry {
			continue
		}

		asset.BumpRelease = asset.Release + 1
		if err := resource.Flush(asset); err != nil {
			log.WithField("file", fname).WithError(err).Errorln("Unable to flush file content")
		}
		log.WithField("file", fname).Println("File correctly updated")
	}
}
