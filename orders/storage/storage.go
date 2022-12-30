package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"
)

var ErrBadData = errors.New("bad data")

type Storage struct {
	db    *Database
	cache map[string][]byte
	sync.RWMutex
}

func New(databaseURL string) (*Storage, error) {
	db, err := NewDatabase(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("storage new: %w", err)
	}

	s := &Storage{
		db:    db,
		cache: make(map[string][]byte),
	}

	return s, nil
}

func (s *Storage) List() ([]string, error) {
	ret, err := s.db.LoadList()
	if err != nil {
		return nil, fmt.Errorf("storage list: %w", err)
	}

	return ret, nil
}

func (s *Storage) Add(orderJson []byte) error {
	var order Order
	if json.Unmarshal(orderJson, &order) != nil {
		return ErrBadData
	}

	if err := s.db.Store(order); err != nil {
		return fmt.Errorf("storage add store: %w", err)
	}

	s.Lock()
	defer s.Unlock()

	s.cache[order.Id] = orderJson
	return nil
}

func (s *Storage) tryGet(id string) ([]byte, bool) {
	s.RLock()
	defer s.RUnlock()

	orderJson, ok := s.cache[id]

	return orderJson, ok
}

func (s *Storage) Get(id string) ([]byte, error) {
	if orderJson, ok := s.tryGet(id); ok {
		return orderJson, nil
	}

	order, err := s.db.Load(id)
	if err != nil {
		return nil, fmt.Errorf("storage get load: %w", err)
	}

	if order == nil {
		return nil, nil
	}

	orderJson, err := json.MarshalIndent(order, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("storage get marshal: %w", err)
	}

	s.Lock()
	defer s.Unlock()

	s.cache[id] = orderJson
	return orderJson, nil
}

func (s *Storage) Close() {
	s.db.Close()
}
