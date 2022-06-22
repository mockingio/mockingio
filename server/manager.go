package server

var _manager *manager

func init() {
	_manager = &manager{
		states: map[string]*State{},
	}
}

func addServer(id string, url string, shutdown func()) {
	_manager.states[id] = &State{
		ServerURL: url,
		Shutdown:  shutdown,
	}
}

type manager struct {
	states map[string]*State
}

func RemoveServer(id string) {
	if state, ok := _manager.states[id]; ok {
		state.Shutdown()
		delete(_manager.states, id)
	}
}

func RemoveAllServers() {
	for _, state := range _manager.states {
		state.Shutdown()
	}
	_manager.states = map[string]*State{}
}

func GetServerURLs() []string {
	var urls []string
	for _, state := range _manager.states {
		urls = append(urls, state.ServerURL)
	}

	return urls
}

type State struct {
	ServerURL string
	Shutdown  func()
}
