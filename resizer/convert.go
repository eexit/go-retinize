package resizer

import (
	"os/exec"
	"fmt"
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

func (ip *convert) WithWidth(img string, w int) error {
	cmd := exec.Command(fmt.Sprintf("%s %s -resize %d %s", ip.bin, img, w, img))

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

