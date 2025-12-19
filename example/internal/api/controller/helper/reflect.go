package helper

import (
	"reflect"
	"strings"
)

// SetupControllers
// auto setup service's controllers by reflect, the controller name must be ended with 'Controller'.
func SetupControllers(service, data any) {
	v := reflect.ValueOf(service)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	dep := reflect.ValueOf(data)

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if strings.HasSuffix(field.Type().Name(), "Controller") {
			for j := 0; j < field.NumField(); j++ {
				ff := field.Field(j)
				if ff.Type() == dep.Type() {
					ff.Set(dep)
				}
			}
		}
	}
}
