package bufqueue

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	ErrReadNoMessages = errors.New("bufqueue: read no messages")
	ErrReadTimeout    = errors.New("bufqueue: read timeout")
)

type blockFileReader struct {
	rc ReaderConfig

	filename string
	fs       *fileSource
	fngBad   *fileNameGenerator
}

func newBlockFileReader(rc ReaderConfig) (*blockFileReader, error) {
	return &blockFileReader{
		rc: rc,
		fs: newFileSource(rc.Dirname, rc.FilePrefix, rc.ReadDuration),
	}, nil
}

func (bfr *blockFileReader) NextFile() (fi os.FileInfo, err error) {

	bfr.RemoveFile()

	fi, err = bfr.fs.NextFile()
	if err != nil {
		return nil, err
	}

	bfr.filename = filepath.Join(bfr.rc.Dirname, fi.Name())

	return fi, nil
}

func (bfr *blockFileReader) Read(block []byte) (int, error) {

	file, err := os.Open(bfr.filename)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	return file.Read(block)
}

func (bfr *blockFileReader) RemoveFile() (err error) {
	if bfr.filename != "" {
		err = os.Remove(bfr.filename)
		bfr.filename = ""
	}
	return
}

func (bfr *blockFileReader) Write(data []byte) error {
	return ioutil.WriteFile(bfr.filename, data, 0600)
}

func (bfr *blockFileReader) MoveFileToBad() error {

	filename := bfr.filename
	if filename == "" {
		return nil
	}

	bfr.filename = ""

	if bfr.fngBad == nil {
		dirname := filepath.Join(bfr.rc.Dirname, "bad")
		fng, err := newFileNameGenerator(dirname, bfr.rc.FilePrefix)
		if err != nil {
			return err
		}
		bfr.fngBad = fng
	}

	badFilename := bfr.fngBad.Generate()

	err := os.Rename(filename, badFilename)
	if err != nil {
		return err
	}

	log.Printf("bad block-file [%s] rename to [%s]", filename, badFilename)

	return nil
}

type fileSource struct {
	dirname      string
	prefix       string
	readDuration time.Duration

	fileList []os.FileInfo
}

func newFileSource(dirname, prefix string, readDuration time.Duration) *fileSource {
	return &fileSource{
		dirname:      dirname,
		prefix:       prefix,
		readDuration: readDuration,
	}
}

func (fs *fileSource) NextFile() (os.FileInfo, error) {

	if len(fs.fileList) == 0 {
		fileList, err := getFilesTimeout(fs.dirname, fs.prefix, fs.readDuration)
		if err != nil {
			return nil, err
		}
		fs.fileList = fileList
	}

	if len(fs.fileList) > 0 {
		fi := fs.fileList[0]
		fs.fileList = fs.fileList[1:]
		return fi, nil
	}

	return nil, ErrReadNoMessages
}

func getFilesTimeout(dirname, prefix string, d time.Duration) (fileList []os.FileInfo, err error) {
	var (
		start    = time.Now()
		sleepDur = 10 * time.Millisecond
	)
	for {
		fileList, err = getFiles(dirname, prefix)
		if err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}
		} else {
			if len(fileList) > 0 {
				return fileList, nil
			}
		}

		t := time.Since(start)
		if t > d {
			return nil, ErrReadTimeout
		}

		if maxDur := d - t; sleepDur > maxDur {
			sleepDur = maxDur
		}
		time.Sleep(sleepDur)
		sleepDur *= 2
	}
}

func getFiles(dirname, prefix string) (fileList []os.FileInfo, err error) {
	fis, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		name := fi.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		if strings.HasSuffix(name, tempExt) {
			continue
		}
		fileList = append(fileList, fi)
	}
	return fileList, nil
}
