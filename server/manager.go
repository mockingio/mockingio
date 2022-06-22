package server

var _manager *manager

func init() {
	_manager = &manager{
		controllers: map[string]*Controller{},
	}
}

type manager struct {
	controllers map[string]*Controller
}

func RemoveServer(id string) {
	if controller, ok := _manager.controllers[id]; ok {
		controller.Shutdown()
		delete(_manager.controllers, id)
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
	for _, controller := range _manager.controllers {
		urls = append(urls, controller.URL)
	}

	return urls
}

func addServer(id string, controller *Controller) {
	_manager.controllers[id] = controller
}
