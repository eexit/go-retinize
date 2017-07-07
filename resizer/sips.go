package resizer

import (
	"fmt"
	"os/exec"
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

func (ip *sips) WithWidth(img string, w int) error {
	cmd := exec.Command(fmt.Sprintf("%s --resampleWidth=%d %s", ip.bin, w, img))

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
