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

	flag.Parse()
}

func main() {
	for _, fname := range flag.Args() {
		log.WithField("file", fname).Println("Inspecting resource")
		asset, err := resource.Parse(fname)
		if err != nil {
			log.WithField("file", fname).WithError(err).Warnln("Unable to parse asset")
		}

		for entry := range asset.Source {
			for source := range asset.Source[entry] {
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
				bump, err := prov.Bump(source)
				if err != nil {
					log.WithField("source", source).WithError(err).Warnln("Unable to bump source")
					continue
				}

				log.WithField("source", bump).Println("Source correctly bumped")
			}
		}
	}
}
