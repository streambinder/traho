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

		// iterate over package sources
		for entry := range asset.Source {
			for source, hash := range asset.Source[entry] {
				log.WithField("source", asset.SourceID(source)).Debugln("Analyzing asset")

				// pick right provider for source
				prov, err := provider.For(source)
				if err != nil {
					log.WithField("source", asset.SourceID(source)).WithError(err).Warnln("Unable to find matching provider")
					continue
				}

				// check source on provider for updates
				log.WithFields(logrus.Fields{
					"provider": prov.Name(),
					"source":   asset.SourceID(source),
				}).Debugln("Bumping source")
				bump, version, err := prov.Bump(source)
				if err != nil {
					log.WithField("source", asset.SourceID(source)).WithError(err).Errorln("Unable to bump source")
					continue
				}

				// use only first source to update version
				if asset.Version == version {
					asset.BumpSource = append(asset.BumpSource, map[string]string{source: hash})
					continue
				}

				// do not fetch data if in dry mode
				if argDry {
					log.WithFields(logrus.Fields{
						"source":  asset.SourceID(source),
						"version": version,
					}).Println("Source has an update")
					continue
				}

				// fetch hash sum of updated source
				log.WithField("source", asset.SourceID(source)).Debugln("Going to fetch source and calculate SHA265 hash")
				assetHash, err := resource.Hash(bump)
				if err != nil {
					log.WithField("source", asset.SourceID(source)).WithError(err).Errorln("Unable to calculate hash for source")
					continue
				}
				asset.BumpSource = append(asset.BumpSource, map[string]string{bump: assetHash})

				// set new package version based
				// on first source update met
				if entry == 0 {
					asset.BumpVersion = version
				}

				log.WithFields(logrus.Fields{
					"source":  asset.SourceID(source),
					"version": version,
				}).Debugln("Source correctly bumped")
			}
		}

		// do not persist update if in dry mode
		if argDry {
			continue
		}

		// check asset has been updated
		if len(asset.BumpVersion) == 0 {
			log.WithField("name", asset.Name).Println("Package not update")
			continue
		}

		// persist update on package file
		asset.BumpRelease = asset.Release + 1
		if err := resource.Flush(asset); err != nil {
			log.WithField("file", fname).WithError(err).Errorln("Unable to flush file content")
		}
		log.WithFields(logrus.Fields{
			"name":    asset.Name,
			"version": asset.BumpVersion,
		}).Println("Package updated")
	}
}
