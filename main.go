package main

import (
	"errors"
	"fmt"
	"github.com/eexit/go-retinize/retinize"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"gopkg.in/urfave/cli.v2"
	"os"
	"reflect"
)

func main() {
	params := retinize.Params{}
	var logLevel, logFormat string

	app := cli.NewApp()
	app.Name = "retinize"
	app.Version = retinize.VERSION
	app.Usage = "Automate down-scaling of images to target lower resolution on non-retina screens"
	app.UsageText = "retinize [global options...] arguments..."
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "factor",
			Usage:       "Highest retinal scaling factor (e.g. 2 -> @2x, 3 -> @3x)",
			Destination: &params.Factor,
		},
		cli.IntFlag{
			Name:        "width",
			Usage:       "Non-retina base image width. Up-scaled images will have {factor} * {base-width} width",
			Destination: &params.Width,
		},
		cli.BoolFlag{
			Name:        "keep-src",
			Usage:       "Makes a copy to the original file instead of resizing (default = false)",
			Destination: &params.BackupSrc,
		},
		cli.StringFlag{
			Name:        "log-level",
			Usage:       "Set log verbosity",
			Value:       "debug",
			Destination: &logLevel,
		},
		cli.StringFlag{
			Name:        "log-format",
			Usage:       "Set log output format (text or json)",
			Value:       "text",
			Destination: &logFormat,
		},
	}
	app.Action = func(ctx *cli.Context) error {
		if ctx.Int("width") == 0 {
			return cli.NewExitError("Missing --width flag", 1)
		}

		if ctx.Int("width") < 0 {
			return cli.NewExitError("Invalid --width flag value", 2)
		}

		err := configureLogger(logLevel, logFormat)
		if err != nil {
			return cli.NewExitError(err, 1)
		}

		retinize.Retinize(params, ctx.Args())

		return nil
	}

	app.Run(os.Args)
}

func configureLogger(level, format string) error {
	ll, err := logrus.ParseLevel(level)
	if err != nil {
		return errors.New(fmt.Sprintf("Invalid log level \"%s\"", level))
	}

	// Sets the log level before the formatter because the formatter might add some
	// extra options according to the log level
	logrus.SetLevel(ll)

	formatter, err := resolveFormatter(format)
	if err != nil {
		return err
	}

	logrus.SetFormatter(formatter)

	logrus.WithFields(logrus.Fields{
		"prefix":    "main",
		"level":     ll,
		"formatter": reflect.Indirect(reflect.ValueOf(formatter)).Type().Name(),
	}).Info("Logger configuration")

	return nil
}

func resolveFormatter(format string) (logrus.Formatter, error) {
	switch format {
	case "text":
		formatter := new(prefixed.TextFormatter)

		if logrus.GetLevel() == logrus.DebugLevel {
			formatter.FullTimestamp = true
		}
		return formatter, nil
	case "json":
		return new(logrus.JSONFormatter), nil
	default:
		return nil, errors.New(fmt.Sprintf("Not suppported logger format: \"%s\"", format))
	}
}
