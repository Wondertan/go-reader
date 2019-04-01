package namer

import (
	"io"
)

type ReadNamer interface {
	io.Reader

	Name() string
}
