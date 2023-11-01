package filedump

import (
	"io"
	"os"
)

type fileReadWriteGetter struct {
	filePath string
}

func New(filePath string) *fileReadWriteGetter {
	return &fileReadWriteGetter{
		filePath: filePath,
	}
}

func (f *fileReadWriteGetter) Get() (io.ReadWriteCloser, error) {
	file, err := os.OpenFile(f.filePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	return file, nil
}
