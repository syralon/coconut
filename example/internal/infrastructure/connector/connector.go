package connector

import (
	"github.com/syralon/coconut/example/internal/config"
	"github.com/syralon/coconut/example/proto/syralon/example"
	"github.com/syralon/coconut/transport/grpc"
)

// Connector holds all the grpc service clients
type Connector struct {
	BookService      example.BookServiceClient
	BookShelfService example.BookShelfServiceClient
}

func NewConnector(conf *config.Config) (*Connector, func(), error) {
	connector := &Connector{}
	cc := grpc.NewConnector()
	err := cc.Connects(
		conf.Connector.BookService.ConnectFunc(&connector.BookService, example.NewBookServiceClient),
		conf.Connector.BookShelfService.ConnectFunc(&connector.BookShelfService, example.NewBookShelfServiceClient),
	)
	if err != nil {
		return nil, nil, err
	}
	return connector, func() { _ = cc.Close() }, err
}
