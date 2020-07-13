package bufq

import (
	"sync"
)

type syncError struct {
	sync.Mutex

	err error
}

func (p *syncError) Store(err error) {
	p.Lock()
	defer p.Unlock()
	if p.err != nil {
		return
	}
	p.err = err
}

func (p *syncError) Load() error {
	p.Lock()
	defer p.Unlock()
	return p.err
}
