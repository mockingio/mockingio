package memory

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/smockyio/smocky/backend/mock/config"
	"github.com/smockyio/smocky/backend/persistent"
)

var _ persistent.Persistent = &Memory{}

type Memory struct {
	mu      sync.Mutex
	configs map[string]*config.Config
	kv      map[string]any
}

func New() *Memory {
	return &Memory{
		configs: map[string]*config.Config{},
		kv:      map[string]any{},
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

func (m *Memory) SetConfig(ctx context.Context, cfg *config.Config) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.configs[cfg.ID] = cfg
	return nil
}

func (m *Memory) GetConfig(ctx context.Context, id string) (*config.Config, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	cfg, ok := m.configs[id]
	if !ok {
		return nil, nil
	}

	return cfg, nil
}

func (m *Memory) GetConfigs(ctx context.Context) ([]*config.Config, error) {
	var configs []*config.Config
	for _, cfg := range m.configs {
		configs = append(configs, cfg)
	}

	return configs, nil
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

func (m *Memory) Increment(_ context.Context, key string) (int, error) {
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

func (m *Memory) SetActiveSession(ctx context.Context, configID string, sessionID string) error {
	return m.Set(ctx, toActiveSessionKey(configID), sessionID)
}

func (m *Memory) GetActiveSession(ctx context.Context, configID string) (string, error) {
	value, err := m.Get(ctx, toActiveSessionKey(configID))
	if err != nil {
		return "", err
	}

	if v, ok := value.(string); ok {
		return v, nil
	}

	return "", errors.New("unable to convert to string value")
}

func toActiveSessionKey(configID string) string {
	return fmt.Sprintf("%s-active-session", configID)
}
