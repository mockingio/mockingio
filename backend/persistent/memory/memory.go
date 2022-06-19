package memory

import "github.com/smockyio/smocky/backend/mock/config"

type Memory struct {
	configs map[string]*config.Config
}

func New() *Memory {
	return &Memory{
		configs: map[string]*config.Config{},
	}
}

func (m *Memory) Set(id string, cfg *config.Config) error {
	m.configs[id] = cfg
	return nil
}

func (m *Memory) Get(id string) (*config.Config, error) {
	return m.configs[id], nil
}
