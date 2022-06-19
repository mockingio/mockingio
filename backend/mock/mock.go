package mock

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/backend/mock/config"
	"github.com/smockyio/smocky/backend/mock/matcher"
	"github.com/smockyio/smocky/backend/session"
)

type Mock struct {
	id          string
	mockFetcher mockFetcher
	session     *session.Session
}

func New(id string, mockFetcher mockFetcher) (*Mock, error) {
	return &Mock{
		id:          id,
		mockFetcher: mockFetcher,
		session:     session.New(),
	}, nil
}

func (m *Mock) Match(req *http.Request) *config.Response {
	cfg, err := m.mockFetcher.Get(m.id)
	if err != nil {
		log.WithError(err).Error("loading mock")
		return nil
	}

	for _, route := range cfg.Routes {
		log.Debugf("Matching route: %v", route.Request)
		response, err := matcher.NewRouteMatcher(route, m.session, req).Match()
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

type mockFetcher interface {
	Get(id string) (*config.Config, error)
}
