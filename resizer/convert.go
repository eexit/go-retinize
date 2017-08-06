package resizer

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
	"reflect"
)

type convert struct {
	bin string
}

func NewConvertResizer() Resizer {
	return &convert{}
}

func (ip *convert) IsInstalled() bool {
	path, err := exec.LookPath("convert")

	if err != nil {
		return false
	}

	ip.bin = path
	return true
}

func (ip *convert) ResampleWidth(logger logrus.FieldLogger, srcImgPath, destImgPath string, width int) bool {
	logger.WithFields(logrus.Fields{
		"prefix": "image-processor",
		"ip":     reflect.Indirect(reflect.ValueOf(ip)).Type().Name(),
	})
	cmdl := fmt.Sprintf("%s --out %s --resampleWidth %d %s", ip.bin, destImgPath, width, srcImgPath)

	logger.WithField("command", cmdl).Debug("Command issue")
	cmd := exec.Command(fmt.Sprintf("%s %s -resize %d %s", ip.bin, srcImgPath, width, destImgPath))

	if err := cmd.Run(); err != nil {
		logger.WithError(err).Error("Command error")
		return false
	}

	return true
}
