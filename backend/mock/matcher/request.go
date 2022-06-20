package matcher

import (
	"context"
	"fmt"
	"net/http"
)

type Request struct {
	HTTPRequest *http.Request
	Session     session
	SessionID   string
}

func (r Request) CountID() string {
	return fmt.Sprintf("%s-%s-%s-count", r.HTTPRequest.Method, r.HTTPRequest.URL, r.SessionID)
}

func (r Request) SequenceID() string {
	return fmt.Sprintf("%s-%s-%s-sequence", r.HTTPRequest.Method, r.HTTPRequest.URL, r.SessionID)
}

type session interface {
	Set(_ context.Context, key string, value any) error
	GetInt(ctx context.Context, key string) (int, error)
	Increase(_ context.Context, key string) (int, error)
}
