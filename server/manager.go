package server

import "errors"

var _manager *manager

func init() {
	_manager = &manager{
		controllers: map[string]*Controller{},
		states:      map[string]State{},
	}
}

type manager struct {
	controllers map[string]*Controller
	states      map[string]State
}

func GetStates() map[string]State {
	return _manager.states
}

func InitState(mockID string) {
	_manager.states[mockID] = State{
		MockID: mockID,
		URL:    "",
		Status: Stopped,
	}
}

func GetState(mockID string) State {
	return _manager.states[mockID]
}

func SetState(mockID, url string, status string) {
	_manager.states[mockID] = State{
		MockID: mockID,
		URL:    url,
		Status: status,
	}
}

func RemoveServer(id string) (State, error) {
	controller, ok := _manager.controllers[id]
	if !ok {
		return State{}, errors.New("mock not found")
	}

	controller.Shutdown()
	delete(_manager.controllers, id)
	InitState(id)

	return _manager.states[id], nil
}

func RemoveAllServers() {
	for _, state := range _manager.controllers {
		state.Shutdown()
	}
	_manager.controllers = map[string]*Controller{}
}

func GetServerURLs() []string {
	var urls []string
	for _, state := range _manager.states {
		if state.Status == Running {
			urls = append(urls, state.URL)
		}
	}

	return urls
}

func addServer(id string, controller *Controller) {
	_manager.controllers[id] = controller
}
