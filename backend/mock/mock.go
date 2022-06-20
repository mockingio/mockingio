package mock

import (
	"context"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/backend/mock/config"
	"github.com/smockyio/smocky/backend/mock/matcher"
)

type Mock struct {
	configID   string
	sessionID  string
	persistent persistent
}

func New(configID, sessionID string, persistent persistent) (*Mock, error) {
	return &Mock{
		configID:   configID,
		sessionID:  sessionID,
		persistent: persistent,
	}, nil
}

func (m *Mock) Match(req *http.Request) *config.Response {
	cfg, err := m.persistent.GetConfig(req.Context(), m.configID)
	if err != nil {
		log.WithError(err).Error("loading mock")
		return nil
	}

	for _, route := range cfg.Routes {
		log.Debugf("Matching route: %v", route.Request)
		response, err := matcher.NewRouteMatcher(route, matcher.Request{
			HTTPRequest: req,
			SessionID:   m.sessionID,
			Session:     m.persistent,
		}).Match()
		if err != nil {
			log.WithError(err).Error("error while matching route")
			continue
		}

		if response == nil {
			log.Debug("no route matched")
			continue
		}

		if response.Delay > 0 {
			time.Sleep(time.Millisecond * time.Duration(response.Delay))
		}

		return response
	}

	return nil
}

func (m *Mock) Handler(w http.ResponseWriter, r *http.Request) {
	response := m.Match(r)
	if response == nil {
		// TODO: no matched? What will be the response?
		w.WriteHeader(http.StatusNotFound)
		return
	}

	for k, v := range response.Headers {
		w.Header().Add(k, v)
	}

	if response.Status == 0 {
		response.Status = 200
	}
	w.WriteHeader(response.Status)
	_, _ = w.Write([]byte(response.Body))
}

type persistent interface {
	GetConfig(ctx context.Context, id string) (*config.Config, error)
	Set(_ context.Context, key string, value any) error
	GetInt(ctx context.Context, key string) (int, error)
	Increase(_ context.Context, key string) (int, error)
}
