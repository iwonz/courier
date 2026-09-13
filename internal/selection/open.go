package selection

import (
	"io"
	"os"
)

func defaultOpenFile(name string) (io.ReadCloser, error) { return os.Open(name) }
