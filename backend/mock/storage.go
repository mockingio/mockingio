package mock

import (
	"context"

	"github.com/smockyio/smocky/backend/session"
)

type storage interface {
	Save(ctx context.Context, session *session.Session) error
	Load(ctx context.Context, sessionID string) (*session.Session, error)
}
