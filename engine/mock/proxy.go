package mock

type Proxy struct {
	Enabled         bool              `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Host            string            `yaml:"host,omitempty" json:"host,omitempty"`
	RequestHeaders  map[string]string `yaml:"request_headers,omitempty" json:"request_headers,omitempty"`
	ResponseHeaders map[string]string `yaml:"response_headers,omitempty" json:"response_headers,omitempty"`
}
