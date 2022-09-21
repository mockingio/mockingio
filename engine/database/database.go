package database

import (
	"context"

	"github.com/mockingio/mockingio/engine/mock"
)

// EngineDB represents the database interface for the engine
type EngineDB interface {
	MockReadWriter
	GetInt(ctx context.Context, mockID, key string) (int, error)
	Increment(ctx context.Context, mockID, key string) (int, error)
	Set(ctx context.Context, mockID, key, value string) error
	Get(ctx context.Context, mockID, key string) (string, error)
	SetActiveSession(ctx context.Context, mockID string, sessionID string) error
	GetActiveSession(ctx context.Context, mockID string) (string, error)
}

type MockReadWriter interface {
	GetMock(ctx context.Context, id string) (*mock.Mock, error)
	SetMock(ctx context.Context, cfg *mock.Mock) error
}

type Database interface {
	EngineDB
	CRUD
}

// CRUD represents the database interface for the CRUD operations
type CRUD interface {
	MockReadWriter
	GetMocks(ctx context.Context) ([]*mock.Mock, error)
	PatchRoute(ctx context.Context, mockID string, routeID string, data string) error
	DeleteRoute(ctx context.Context, mockID string, routeID string) error
	CreateRoute(ctx context.Context, mockID string, data string) error
	PatchResponse(ctx context.Context, mockID, routeID, responseID, data string) error
}
