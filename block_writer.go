package bufqueue

import (
	"io/ioutil"
	"os"
)

const tempExt = ".tmp"

type blockFileWriter struct {
	fng *fileNameGenerator
}

func newBlockFileWriter(fc FilesConfig) (*blockFileWriter, error) {
	fng, err := newFileNameGenerator(fc.Dirname, fc.FilePrefix)
	if err != nil {
		return nil, err
	}
	return &blockFileWriter{fng}, nil
}

func (bfw *blockFileWriter) SaveToFile(block []byte) error {
	filename := bfw.fng.Generate()
	tmp := filename + tempExt
	err := ioutil.WriteFile(tmp, block, 0600)
	if err != nil {
		return err
	}
	return os.Rename(tmp, filename)
}
