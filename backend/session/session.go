package session

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
)

const (
	keyRequestNumber     = "request_number"
	keyNextResponseIndex = "next_response_index"
)

type routeState map[string]any

type Session struct {
	mu          sync.Mutex
	routeStates map[string]routeState
}

func New() *Session {
	return &Session{
		routeStates: map[string]routeState{},
	}
}

func (s *Session) Get(method, path, key string) any {
	s.mu.Lock()
	defer s.mu.Unlock()

	r, ok := s.routeStates[RouteID(method, path)]
	if !ok {
		return nil
	}

	return r[key]
}

func (s *Session) Set(method, path, key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := RouteID(method, path)

	r, ok := s.routeStates[id]
	if !ok {
		r = routeState{}
		s.routeStates[id] = r
	}

	r[key] = value
	s.routeStates[id] = r
}

func (s *Session) GetInt(method, path, key string) int {
	v := s.Get(method, path, key)
	if v == nil {
		return 0
	}

	value, ok := v.(int)
	if !ok {
		return 0
	}

	return value
}

func (s *Session) Increase(method, path, key string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := RouteID(method, path)

	r, ok := s.routeStates[id]
	if !ok {
		r = routeState{}
		r[key] = 1
		s.routeStates[id] = r

		return 1
	}

	val, ok := r[key].(int)
	if !ok {
		return 0
	}

	val++
	r[key] = val
	s.routeStates[RouteID(method, path)] = r

	return val
}

func RouteID(method, path string) string {
	method = strings.ToLower(method)
	return fmt.Sprintf("%v:%v", method, path)
}

func (s *Session) NextResponseIndex(req *http.Request) int {
	return s.GetInt(req.Method, req.URL.Path, keyNextResponseIndex)
}

func (s *Session) SetNextResponseIndex(req *http.Request, idx int) {
	s.Set(req.Method, req.URL.Path, keyNextResponseIndex, idx)
}

func (s *Session) GetRequestNumber(req *http.Request) int {
	return s.GetInt(req.Method, req.URL.Path, keyRequestNumber)
}

func (s *Session) IncreaseRequestNumber(req *http.Request) int {
	return s.Increase(req.Method, req.URL.Path, keyRequestNumber)
}
