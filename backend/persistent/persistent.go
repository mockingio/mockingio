package persistent

import (
	"context"

	"github.com/smockyio/smocky/backend/mock/config"
)

type Persistent interface {
	SetConfig(ctx context.Context, id string, cfg *config.Config) error
	GetConfig(ctx context.Context, id string) (*config.Config, error)
	Set(_ context.Context, key string, value any) error
	Get(_ context.Context, key string) (any, error)
	GetInt(ctx context.Context, key string) (int, error)
	Increase(_ context.Context, key string) (int, error)
}
