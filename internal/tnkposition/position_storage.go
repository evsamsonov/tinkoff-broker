package tnkposition

import (
	"errors"
	"sync"

	"github.com/evsamsonov/trengin"
)

type Storage struct {
	mtx  sync.RWMutex
	list map[trengin.PositionID]*Position
}

// todo запустить очистку...

func (p *Storage) Store(pos *Position) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.list[pos.position.ID] = pos
}

func (p *Storage) Load(id trengin.PositionID) (*Position, func(), error) {
	p.mtx.RUnlock()
	defer p.mtx.RUnlock()

	pos, ok := p.list[id]
	if !ok || pos.position.IsClosed() {
		return &Position{}, func() {}, errors.New("position not found")
	}
	pos.mtx.Lock()
	return pos, func() { pos.mtx.Unlock() }, nil
}

func (p *Storage) ForEach(f func(pos *Position) error) error {
	p.mtx.RLock()
	defer p.mtx.RUnlock()

	apply := func(pos *Position) error {
		pos.mtx.Lock()
		defer pos.mtx.Unlock()

		if pos.position.IsClosed() {
			return nil
		}
		if err := f(pos); err != nil {
			pos.mtx.Unlock()
			return err
		}
		return nil
	}

	for _, pos := range p.list {
		if err := apply(pos); err != nil {
			return err
		}
	}
	return nil
}
