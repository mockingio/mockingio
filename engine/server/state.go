package server

const (
	Running = "running"
	Stopped = "stopped"
)

type MockServerState struct {
	MockID           string `json:"mock_id"`
	URL              string `json:"url"`
	Status           string `json:"status"`
	shutdownServerFn func()
}

func (s *MockServerState) shutdownServer() {
	if s.shutdownServerFn != nil {
		s.shutdownServerFn()
	}
	s.Status = Stopped
	s.URL = ""
}
