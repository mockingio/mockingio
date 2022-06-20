package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/smockyio/smocky/backend/mock/config"
)

type Memory struct {
	mu sync.Mutex
	kv map[string]any
}

func New() *Memory {
	return &Memory{
		kv: map[string]any{},
	}
}

func (m *Memory) Get(_ context.Context, key string) (any, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.kv[key], nil
}

func (m *Memory) Set(_ context.Context, key string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.kv[key] = value
	return nil
}

func (m *Memory) SetConfig(ctx context.Context, id string, cfg *config.Config) error {
	return m.Set(ctx, id, cfg)
}

func (m *Memory) GetConfig(ctx context.Context, id string) (*config.Config, error) {
	value, err := m.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if value == nil {
		return nil, nil
	}

	cfg, ok := value.(*config.Config)
	if !ok {
		return nil, nil
	}

	return cfg, nil
}

func (m *Memory) GetInt(ctx context.Context, key string) (int, error) {
	v, err := m.Get(ctx, key)
	if err != nil {
		return 0, err
	}

	if v == nil {
		return 0, nil
	}

	value, ok := v.(int)
	if !ok {
		return 0, nil
	}

	return value, nil
}

func (m *Memory) Increase(_ context.Context, key string) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	value, ok := m.kv[key]
	if !ok {
		m.kv[key] = 1
		return 1, nil
	}

	val, ok := value.(int)
	if !ok {
		return 0, errors.New(fmt.Sprintf("unable to increase non-int key (%s)", key))
	}

	val++
	m.kv[key] = val

	return val, nil
}
