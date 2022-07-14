package server

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

func SetState(mockId, url string, status string) {
	_manager.states[mockId] = State{
		MockID: mockId,
		URL:    url,
		Status: status,
	}
}

func RemoveServer(id string) {
	if controller, ok := _manager.controllers[id]; ok {
		controller.Shutdown()
		delete(_manager.controllers, id)
		InitState(id)
	}
}

func PauseServer(id string) {
	if controller, ok := _manager.controllers[id]; ok {
		controller.Pause()
	}
}

func ResumeServer(id string) {
	if controller, ok := _manager.controllers[id]; ok {
		controller.Resume()
	}
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
