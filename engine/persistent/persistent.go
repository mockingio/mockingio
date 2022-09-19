package persistent

import (
	"context"

	"github.com/mockingio/mockingio/engine/mock"
)

type Engine interface {
	GetMock(ctx context.Context, id string) (*mock.Mock, error)
	GetInt(ctx context.Context, key string) (int, error)
	Increment(ctx context.Context, key string) (int, error)
	Set(ctx context.Context, key string, value any) error
	Get(ctx context.Context, key string) (any, error)
	SetActiveSession(ctx context.Context, mockID string, sessionID string) error
	GetActiveSession(ctx context.Context, mockID string) (string, error)
}

type Persistent interface {
	Engine
	API
}

type API interface {
	GetMock(ctx context.Context, id string) (*mock.Mock, error)
	SetMock(ctx context.Context, cfg *mock.Mock) error
	GetMocks(ctx context.Context) ([]*mock.Mock, error)
	PatchRoute(ctx context.Context, mockID string, routeID string, data string) error
	DeleteRoute(ctx context.Context, mockID string, routeID string) error
	CreateRoute(ctx context.Context, mockID string, data string) error
	PatchResponse(ctx context.Context, mockID, routeID, responseID, data string) error
}
