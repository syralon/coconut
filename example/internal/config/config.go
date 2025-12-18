package config

import (
	"github.com/syralon/coconut/pkg/etcdutil"
	"github.com/syralon/coconut/transport/gateway"
	"github.com/syralon/coconut/transport/grpc"
)

type Config struct {
	GRPC    grpc.Config     `json:"grpc"    yaml:"grpc"`
	Gateway gateway.Config  `json:"gateway" yaml:"gateway"`
	ETCD    etcdutil.Config `json:"etcd" yaml:"etcd"`
}
