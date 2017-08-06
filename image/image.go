package image

import (
	"github.com/sirupsen/logrus"
	goimg "image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	FileInfo  FileInfo
	ImageInfo ImageInfo
}

type FileInfo struct {
	Path, Dir, Name, Ext string
}

type ImageInfo struct {
	Type  string
	Width int
}

func Parse(path string) *Image {
	logger := logrus.WithFields(logrus.Fields{
		"prefix": "image-parser",
		"file":   path,
	})

	logger.Debug("Stat-ing file")

	_, err := os.Stat(path)

	if err != nil {
		if os.IsNotExist(err) {
			logger.WithError(err).Warn("File not found")
			return nil
		}

		if os.IsPermission(err) {
			logger.WithError(err).Warn("Permission denied")
			return nil
		}

		logger.Error(err)
		return nil
	}

	logger.Debug("Opening file")
	file, err := os.Open(path)

	if err != nil {
		if os.IsPermission(err) {
			logger.WithError(err).Warn("Permission denied")
			return nil
		}

		logger.Error(err)
		return nil
	}

	defer file.Close()

	fileInfo := FileInfo{
		Path: path,
		Dir:  filepath.Dir(path),
		Name: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		Ext:  filepath.Ext(path),
	}

	logger.Debug("Decoding file")
	i, format, err := goimg.DecodeConfig(file)

	if err != nil {
		if err == goimg.ErrFormat {
			logger.WithError(err).Warn("Unsupported file type")
			return nil
		}

		logger.Error(err)
		return nil
	}

	imgInfo := ImageInfo{
		Type:  format,
		Width: i.Width,
	}

	img := Image{
		FileInfo:  fileInfo,
		ImageInfo: imgInfo,
	}

	logger.WithField("image", img).Debug("Parsed image")

	return &img
}
