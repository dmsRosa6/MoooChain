package store

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
)

var ErrKeyNotFound = errors.New("key not found")

type MemStore struct {
    mu   sync.RWMutex
    data map[string]string
}

func NewMemStore() *MemStore {
    return &MemStore{
        data: make(map[string]string),
    }
}

func (m *MemStore) Get(ctx context.Context, key string) (string, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()

    val, ok := m.data[key]
    if !ok {
        return "", ErrKeyNotFound
    }
    return val, nil
}

func (m *MemStore) Set(ctx context.Context, key string, value interface{}) error {
    var data string

    switch v := value.(type) {
    case string:
        data = v
    default:
        b, err := json.Marshal(v)
        if err != nil {
            return err
        }
        data = string(b)
    }

    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = data
    return nil
}

