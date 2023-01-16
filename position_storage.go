package tnkbroker

import (
	"errors"
	"sync"

	"github.com/evsamsonov/trengin"
)

type openPositionStorage struct {
	mtx  sync.RWMutex
	list map[trengin.PositionID]*openPosition
}

func (p *openPositionStorage) Store(pos *openPosition) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.list[pos.position.ID] = pos
}

func (p *openPositionStorage) Load(id trengin.PositionID) (*openPosition, func(), error) {
	p.mtx.RUnlock()
	defer p.mtx.RUnlock()

	pos, ok := p.list[id]
	if !ok {
		return &openPosition{}, func() {}, errors.New("position not found")
	}
	pos.mtx.Lock()
	return pos, func() { pos.mtx.Unlock() }, nil
}

func (p *openPositionStorage) ForEach(f func(pos *openPosition) error) error {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	for _, pos := range p.list {
		pos.mtx.Lock()
		if err := f(pos); err != nil {
			pos.mtx.Unlock()
			return err
		}
		pos.mtx.Unlock()
	}
	return nil
}
