package resizer

type Resizer interface {
	IsInstalled() bool
	WithWidth(img string, w int) error
}

