package common

import (
	"github.com/google/uuid"

	"github.com/syralon/coconut/mesh"
	"github.com/syralon/coconut/toolkit/netutil"
	"github.com/syralon/coconut/toolkit/text"
)

type Config struct {
	Name     string        `json:"name"      yaml:"name"`
	ListenOn string        `json:"listen_on" yaml:"listen_on"`
	Timeout  text.Duration `json:"timeout"   yaml:"timeout"`
}

type HTTPConfig struct {
	Config
	TLS TLSConfig `json:"tls" yaml:"tls"`
}

type TLSConfig struct {
	CertFile string `json:"cert_file" yaml:"cert_file"`
	KeyFile  string `json:"key_file"  yaml:"key_file"`
}

func (c *Config) Endpoint() *mesh.Endpoint {
	ip, port := netutil.FingerOut(c.ListenOn, true)
	endpoint := &mesh.Endpoint{
		ID:   uuid.New().String(),
		Name: c.Name,
		Host: ip,
		Port: uint16(port),
	}
	return endpoint
}
