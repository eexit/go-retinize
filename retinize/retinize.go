package retinize

import (
	"fmt"
	"github.com/eexit/go-retinize/image"
	"github.com/eexit/go-retinize/resizer"
	"github.com/sirupsen/logrus"
	"math"
	"os"
	"path/filepath"
)

const VERSION = "0.1.0-dev"
const modifier = "@"

type Params struct {
	Factor, Width int
	BackupSrc     bool
}

type Variant struct {
	ID, Width  int
	Path, Name string
}

func Retinize(params Params, args []string) {
	logrus.WithField("params", params).Info("Retinize parameters")

	r := resizer.ResolveResizer(
		resizer.NewSipsResizer(),
		resizer.NewConvertResizer(),
	)

	for img := range parse(args...) {
		logger := logrus.WithField("image", img.FileInfo.Path)
		var doneVariantNames []string

		if params.BackupSrc {
			if ok := backupSrc(logger, img); !ok {
				continue
			}
		}

		for v := range widthVariants(logger, params.Factor, params.Width, img) {
			if ok := r.ResampleWidth(logger, img.FileInfo.Path, v.Path, v.Width); ok {
				doneVariantNames = append(doneVariantNames, v.Name)
			}
		}

		go func() {
			if len(doneVariantNames) > 0 {
				logger.WithField("done", doneVariantNames).Info("Image processed")
			}
		}()
	}
}

func parse(args ...string) <-chan *image.Image {
	out := make(chan *image.Image)

	go func() {
		for _, arg := range args {
			img := image.Parse(arg)
			if img == nil {
				continue
			}

			out <- img
		}
		close(out)
	}()

	return out
}

func backupSrc(logger logrus.FieldLogger, img *image.Image) bool {
	dest := filepath.Join(img.FileInfo.Dir, fmt.Sprintf("backup-%s%s", img.FileInfo.Name, img.FileInfo.Ext))
	logger = logger.WithField("prefix", "image-backup")
	logger.WithField("backup", dest).Info("Backuping image")

	_, err := os.Stat(dest)

	if err != nil {
		if os.IsExist(err) {
			logger.WithError(err).Warn("Backup image already exist")
			return false
		}
	}

	err = os.Link(img.FileInfo.Path, dest)

	if err != nil {
		logger.WithError(err).Error("Failed to create image backup")
		return false
	}

	return true
}

func widthVariants(logger logrus.FieldLogger, factor, baseWidth int, img *image.Image) <-chan *Variant {
	out := make(chan *Variant)

	logger = logger.WithFields(logrus.Fields{
		"prefix": "variant-computer",
	})

	go func() {
		if img.ImageInfo.Width <= baseWidth {
			logger.Infof("Ignoring image: width is <= %d", baseWidth)
			return
		}

		vc := math.Floor(float64(img.ImageInfo.Width / baseWidth))

		// If retinize --factor param is passed, the loop count will be
		// restrained to the lowest value between available variant count
		// and scaling factor
		if factor > 0 {
			vc = math.Min(float64(factor), vc)
		}

		for i := int(vc); i > 0; i-- {
			path := img.FileInfo.Path
			name := img.FileInfo.Name

			if i > 1 {
				name := fmt.Sprintf("%s%s%d%s", img.FileInfo.Name, modifier, i, img.FileInfo.Ext)
				path = filepath.Join(img.FileInfo.Dir, name)
			}

			v := &Variant{
				ID:    i,
				Width: baseWidth * i,
				Path:  path,
				Name:  name,
			}

			logger.WithField("variant", v).Debug("Computed variant")
			out <- v
		}

		close(out)
	}()

	return out
}
