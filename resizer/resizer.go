package resizer

import "github.com/sirupsen/logrus"

type Resizer interface {
	IsInstalled() bool
	ResampleWidth(logger logrus.FieldLogger, srcImgPath, destImgPath string, width int) bool
}

func ResolveResizer(resizers ...Resizer) Resizer {
	installed := make(chan Resizer)
	defer close(installed)

	for _, rs := range resizers {
		go func(rs Resizer) {
			if rs.IsInstalled() {
				installed <- rs
			}
		}(rs)
	}

	return <-installed
}
