package bufqueue

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gitchander/bufqueue/seqt"
)

type fileNameGenerator struct {
	dirname string
	prefix  string

	tab *seqt.Table
	seq *seqt.Sequence

	val string
}

func newFileNameGenerator(dirname, prefix string) (fng *fileNameGenerator, err error) {

	fng = &fileNameGenerator{
		dirname: dirname,
		prefix:  prefix,

		tab: seqt.NewTable(seqt.UPPER_LETTERS),
		seq: new(seqt.Sequence),
	}

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		if os.IsNotExist(err) {
			if err = os.Mkdir(dirname, 0700); err != nil {
				return nil, err
			}
			return fng, nil // success!
		}
		return nil, err
	}

	var roots []string

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if name := file.Name(); strings.HasPrefix(name, prefix) {
			roots = append(roots, name[len(prefix):])
		}
	}

	sort.Sort(sort.Reverse(sort.StringSlice(roots)))

	for _, root := range roots {
		if seq, err := fng.tab.Parse(root); err == nil {
			fng.seq = seq
			break
		}
	}

	return fng, nil
}

func (fng *fileNameGenerator) Generate() string {
	fng.seq.Next()
	root := fng.tab.String(fng.seq)
	fng.val = filepath.Join(fng.dirname, fng.prefix+root)
	return fng.val
}

func (fng *fileNameGenerator) Value() string {
	return fng.val
}
