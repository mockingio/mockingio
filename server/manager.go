package server

import (
	"fmt"
)

var _manager *manager

func init() {
	_manager = &manager{
		mocks: map[string]item{},
	}
}

type item struct {
	controller *Controller
	state      State
}

type manager struct {
	mocks map[string]item
}

func GetStates() map[string]State {
	states := map[string]State{}
	for id, mock := range _manager.mocks {
		states[id] = mock.state
	}
	return states
}

func NewState(mockID, url string, status string) State {
	return State{
		MockID: mockID,
		URL:    url,
		Status: status,
	}
}

func RemoveServer(id string) (State, error) {
	mock, ok := _manager.mocks[id]
	if !ok {
		return State{}, fmt.Errorf("mock id: %v not found", id)
	}

	mock.controller.Shutdown()
	mock.state = NewState(id, "", Stopped)
	_manager.mocks[id] = mock

	return mock.state, nil
}

func RemoveAllServers() {
	for id, _ := range _manager.mocks {
		_, _ = RemoveServer(id)
	}
}

func GetServerURLs() []string {
	var urls []string
	for _, mock := range _manager.mocks {
		if mock.state.Status == Running {
			urls = append(urls, mock.state.URL)
		}
	}

	return urls
}

func addServer(id string, controller *Controller, state State) {
	_manager.mocks[id] = item{
		controller: controller,
		state:      state,
	}
}
