package bookservice

import (
	"reflect"
	"strings"
	"testing"
)

func TestSetup(t *testing.T) {
	dep := &Dependency{}
	s := &BookService{}
	val := reflect.ValueOf(s).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		t.Log(field.IsValid(), field.Type().Name(), field.Type().String())
		if strings.HasSuffix(field.Type().Name(), "Controller") {
			for j := 0; j < field.NumField(); j++ {
				ff := field.Field(j)
				t.Log(ff.String(), ff.Type() == reflect.TypeOf(dep))
			}
		}
	}
}
