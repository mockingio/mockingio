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
	cfg     *config.Config
	session *session.Session
}

func NewFromYaml(file string) (*Mock, error) {
	cfg, err := config.FromYamlFile(file)
	if err != nil {
		return nil, err
	}

	return New(cfg)
}

func New(cfg *config.Config) (*Mock, error) {
	return &Mock{
		cfg:     cfg,
		session: session.New(),
	}, nil
}

func (m *Mock) Port() string {
	return m.cfg.Port
}

func (m *Mock) Match(req *http.Request) *config.Response {
	for _, route := range m.cfg.Routes {
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
