package mock

type TLS struct {
	Enabled     bool   `yaml:"enabled,omitempty" json:"enabled,omitempty"`
	PEMCertPath string `yaml:"pem_cert_path,omitempty" json:"pem_cert_path,omitempty"`
	PEMKeyPath  string `yaml:"pem_key_path,omitempty" json:"pem_key_path,omitempty"`
}
