package etcdutil

import clientv3 "go.etcd.io/etcd/client/v3"

type Config struct {
	// Endpoints is a list of URLs.
	Endpoints []string `json:"endpoints" yaml:"endpoints"`
	// Username is a username for authentication.
	Username string `json:"username"`
	// Password is a password for authentication.
	Password string `json:"password"`
}

func (c *Config) Config() clientv3.Config {
	return clientv3.Config{
		Endpoints: c.Endpoints,
		Username:  c.Username,
		Password:  c.Password,
	}
}

func (c *Config) NewClient() (*clientv3.Client, error) {
	return clientv3.New(c.Config())
}

func NewClient(c *Config) (*clientv3.Client, error) {
	return c.NewClient()
}
