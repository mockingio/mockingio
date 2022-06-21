package persistent

import (
	"context"

	"github.com/smockyio/smocky/engine/mock"
)

var _default Persistent

func New(p Persistent) {
	_default = p
}

func GetDefault() Persistent {
	if _default == nil {
		panic("default persistent needs to be initialised")
	}
	return _default
}

type Persistent interface {
	SetConfig(ctx context.Context, cfg *mock.Mock) error
	GetConfig(ctx context.Context, id string) (*mock.Mock, error)
	GetConfigs(ctx context.Context) ([]*mock.Mock, error)

	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)

	GetInt(ctx context.Context, key string) (int, error)
	Increment(_ context.Context, key string) (int, error)

	SetActiveSession(ctx context.Context, mockID string, sessionID string) error
	GetActiveSession(ctx context.Context, mockID string) (string, error)
}
