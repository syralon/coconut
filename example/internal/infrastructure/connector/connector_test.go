package connector

import (
	"reflect"
	"testing"

	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/proto/syralon/example"
	cocogrpc "github.com/syralon/coconut/transport/grpc"
)

func TestConnector(t *testing.T) {
	cc, clean, err := NewConnector(&config.Config{Connector: config.ConnectorConfig{
		BookService:      cocogrpc.ServiceClientConfig[example.BookServiceClient]{ClientConfig: cocogrpc.ClientConfig{Target: "127.0.0.1:8000"}},
		BookShelfService: cocogrpc.ServiceClientConfig[example.BookShelfServiceClient]{ClientConfig: cocogrpc.ClientConfig{Target: "127.0.0.1:8000"}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	defer clean()
	t.Log(cc.BookService)
	t.Log(reflect.TypeOf(cc.BookService).String())
}
