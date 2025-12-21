package connector

import (
	"fmt"
	"reflect"

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

func (c *Connector) setup(cfg *config.ConnectorConfig) error {
	cc := grpc.NewConnector()
	connector := reflect.ValueOf(c).Elem()
	confValue := reflect.ValueOf(cfg).Elem()

	for i := 0; i < connector.NumField(); i++ {
		field := connector.Field(i)
		conf := confValue.FieldByName(field.Type().Name())
		if !conf.IsValid() {
			return fmt.Errorf("service %s Config not found", field.Type().Name())
		}
		if conf.Kind() == reflect.Ptr {
			conf = conf.Elem()
		}
		fn := conf.FieldByName("ConnectFunc")
		if !fn.IsValid() {
			return fmt.Errorf("service %s Config has no function 'ConnectFunc'", field.Type().Name())
		}
		fn.Call([]reflect.Value{
			field.Addr(),
		})
	}

}
