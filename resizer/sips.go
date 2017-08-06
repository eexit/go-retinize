package resizer

import (
	"github.com/sirupsen/logrus"
	"os/exec"
	"reflect"
	"strconv"
)

type sips struct {
	bin string
}

func NewSipsResizer() Resizer {
	return &sips{}
}

func (ip *sips) IsInstalled() bool {
	path, err := exec.LookPath("sips")

	if err != nil {
		return false
	}

	ip.bin = path
	return true
}

func (ip *sips) ResampleWidth(logger logrus.FieldLogger, srcImgPath, destImgPath string, width int) bool {
	logger = logger.WithFields(logrus.Fields{
		"prefix": "image-processor",
		"ip":     reflect.Indirect(reflect.ValueOf(ip)).Type().Name(),
	})

	params := []string{"--out", destImgPath, "--resampleWidth", strconv.Itoa(width), srcImgPath}
	logger.WithFields(logrus.Fields{
		"bin":    ip.bin,
		"params": params,
	}).Debug("Command issue")

	cmd := exec.Command(ip.bin, params...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		logger.WithError(err).Error("Command error")
		return false
	}

	logger.WithField("output", string(output)).Debug("Command output")
	return true
}
