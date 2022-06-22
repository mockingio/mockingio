package engine

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/smockyio/smocky/engine/matcher"
	"github.com/smockyio/smocky/engine/mock"
	"github.com/smockyio/smocky/engine/persistent"
)

type Engine struct {
	mockID   string
	isPaused bool
}

func New(mockID string) *Engine {
	return &Engine{
		mockID: mockID,
	}
}

func (eng *Engine) Resume() {
	eng.isPaused = false
}

func (eng *Engine) Pause() {
	eng.isPaused = true
}

func (eng *Engine) Match(req *http.Request) *mock.Response {
	ctx := req.Context()
	db := persistent.GetDefault()
	mok, err := db.GetMock(ctx, eng.mockID)
	if err != nil {
		log.WithError(err).Error("loading mock")
		return nil
	}

	sessionID, err := db.GetActiveSession(ctx, eng.mockID)
	if err != nil {
		log.WithError(err).WithField("config_id", eng.mockID).Error("get active session")
	}

	for _, route := range mok.Routes {
		log.Debugf("Matching route: %v", route.Request)
		response, err := matcher.NewRouteMatcher(route, matcher.Context{
			HTTPRequest: req,
			SessionID:   sessionID,
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

func (eng *Engine) Handler(w http.ResponseWriter, r *http.Request) {
	if eng.isPaused {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	response := eng.Match(r)
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
