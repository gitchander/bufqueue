package bufq

import (
	"container/list"
	"io/ioutil"
	"strings"
	"sync"
)

type FilesConfig struct {
	Dirname    string
	FilePrefix string
}

func ReadFiles(fc FilesConfig) ([]string, error) {
	fis, err := ioutil.ReadDir(fc.Dirname)
	if err != nil {
		// if os.IsNotExist(err) {
		// 	err = os.Mkdir(fc.Dirname, 0755)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	return nil, nil
		// }
		return nil, err
	}
	prefix := fc.FilePrefix
	var files []string
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if name := fi.Name(); strings.HasPrefix(name, prefix) {
			files = append(files, name[len(prefix):])
		}
	}
	return files, nil
}

type fileList struct {
	guard sync.Mutex
	l     *list.List
}

func newFileList() *fileList {
	return &fileList{
		l: list.New(),
	}
}

func (p *fileList) Store(filename string) {
	p.guard.Lock()
	p.l.PushBack(filename)
	p.guard.Unlock()
}

func (p *fileList) Load() string {
	p.guard.Lock()
	e := p.l.Front()
	p.l.Remove(e)
	p.guard.Unlock()
	return e.Value.(string)
}
