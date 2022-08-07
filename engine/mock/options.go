package mock

type mockOptions struct {
	idGeneration bool
}

type Option func(*mockOptions)

func WithIDGeneration() Option {
	return func(m *mockOptions) {
		m.idGeneration = true
	}
}
