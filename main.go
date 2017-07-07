package main

import (
	"os"

	"fmt"
	"github.com/eexit/go-retinize/resizer"
	"gopkg.in/urfave/cli.v2"
)

const VERSION = "0.1.0-dev"

func main() {
	var factor, baseWidth, baseHeight int

	app := cli.NewApp()
	app.Name = "retinize"
	app.Version = VERSION
	app.Usage = "Automate down-scaling of images to target lower resolution or non-retina screens"
	app.UsageText = "retinize [global options...] arguments..."
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:        "factor",
			Usage:       "Highest retinal scaling factor (e.g. 2 -> @2x, 3 -> @3x)",
			Value:       2,
			Destination: &factor,
		},
		cli.IntFlag{
			Name:        "base-width",
			Usage:       "Non-retina base image width. Up-scaled images will have {factor} * {base-width} width",
			Destination: &baseWidth,
		},
		cli.IntFlag{
			Name:        "base-height",
			Usage:       "Non-retina base image height. Up-scaled images will have {factor} * {base-height} height",
			Destination: &baseHeight,
		},
	}
	app.Action = func(c *cli.Context) error {
		fmt.Printf("Images: %s\n", c.Args())
		return retinize()
	}

	app.Run(os.Args)
}

func retinize() error {
	r := resolveResizer(
		resizer.NewSipsResizer(),    // macOS scriptable image processing system
		resizer.NewConvertResizer(), // ImageMagik image converter
	)

	fmt.Printf("%+v\n", r)
	return nil
}

func resolveResizer(resizers ...resizer.Resizer) resizer.Resizer {
	installed := make(chan resizer.Resizer)
	defer close(installed)

	for _, rs := range resizers {
		go func(rs resizer.Resizer) {
			if rs.IsInstalled() {
				installed <- rs
			}
		}(rs)
	}

	return <-installed
}
