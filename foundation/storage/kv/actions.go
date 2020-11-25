package kv

import (
	"github.com/sonemas/libereco/foundation/storage"
)

func (s *KVStorage) Init(state map[string][]byte, _ map[string]interface{}) error {
	if s.version == 0 {
		s.storage = state
	}

	return nil
}

func (s *KVStorage) Set(k string, v []byte, _ map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.storage[k] = v
	s.version++

	return s.save()
}

func (s *KVStorage) Get(k string, _ map[string]interface{}) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	v, ok := s.storage[k]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return v, nil
}

func (s *KVStorage) Fetch(args map[string]interface{}) (map[string][]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.allowFetching {
		return nil, storage.ErrOperationNotAllowed
	}

	// TODO: Implement advanced fetching
	return s.storage, nil
}

func (s *KVStorage) Delete(k string, _ map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.allowDeleting {
		return storage.ErrOperationNotAllowed
	}

	if _, exists := s.storage[k]; !exists {
		return storage.ErrNotFound
	}

	delete(s.storage, k)
	s.version++
	return s.save()
}
